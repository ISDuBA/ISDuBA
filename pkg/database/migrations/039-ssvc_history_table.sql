-- This file is Free Software under the Apache-2.0 License
-- without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
--
-- SPDX-License-Identifier: Apache-2.0
--
-- SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
-- Software-Engineering: 2026 Intevation GmbH <https://intevation.de>


BEGIN;

-- actor taken from event_log. Nice to have, but not sure if necessary? (E.g. if user is deleted, or a users changes have to be rolled back?)
-- changedate taken from event_log, time of change
-- change_number Primary key to differentiate
-- documents_id which document is the subject
-- ssvc current ssvc, taken from event_log of next if exists, from current document otherwise
-- assumes we cannot delete ssvc yet.
CREATE TABLE ssvc_history (
    actor         varchar,
    changedate    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    change_number bigint      NOT NULL DEFAULT 0 ,
    documents_id  integer     NOT NULL                            REFERENCES documents(id) ON DELETE CASCADE,
    ssvc          text

    PRIMARY KEY (documents_id, change_number)
);

-- Gather the data from the old structure
-- fill in all but change_number which is automatic via
INSERT INTO ssvc_history (actor, changedate, change_number, documents_id,  ssvc)
-- construct timeline of events
WITH timeline AS (
-- Take the current values, except for prev_ssvc where the LEAD (next) is taken
-- NULL if last can later be used in COALESCE -> This is basically already what is wanted
    SELECT
        actor,
        time                          AS changedate,
        ROW_NUMBER() OVER (
            PARTITION BY documents_id
            ORDER BY time ASC
        ) AS change_number
        documents_id,
        LEAD(prev_ssvc) OVER (
            PARTITION BY documents_id
            ORDER BY time ASC
        ) AS next_prev_ssvc
    FROM events_log
    WHERE event IN ('add_sscv', 'change_sscv')
)

-- Take what is created in the last step as the timeline, but if the ssvc is empty
-- (no next ssvc) take the current document (d.id) ssvc as ssvc
SELECT
    t.actor,
    t.changedate,
    t.change_number,
    t.documents_id,
    COALESCE(t.next_prev_ssvc, d.ssvc) AS ssvc
FROM timeline t
JOIN documents d ON t.documents_id = d.id;

-- Drop no longer necessary column
ALTER TABLE documents DROP COLUMN IF EXISTS ssvc;
ALTER TABLE events_log DROP COLUMN IF EXISTS prev_ssvc;

COMMIT;


CREATE FUNCTION log_ssvc_history_to_events()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    v_prev_ssvc text;
    v_event     events;
BEGIN
    -- Find the most recent previous SSVC value for this document
    SELECT ssvc
      INTO v_prev_ssvc
      FROM ssvc_history
     WHERE documents_id = NEW.documents_id
       AND change_number < NEW.change_number
     ORDER BY change_number DESC
     LIMIT 1;

    -- Determine event based on the change
    IF v_prev_ssvc IS NULL THEN
        -- No previous value existed: add
        v_event := 'add_sscv';
        -- Not possible yet, but future-proofing:
    ELSIF NEW.ssvc IS NULL THEN
        -- Had value before, now null: delete
        v_event := 'delete_sscv';
    ELSE
        -- Only a change: change
        v_event := 'change_sscv';
    END IF;

    INSERT INTO events_log (
        event,
        state,
        time,
        actor,
        documents_id,
        comments_id
    )
    VALUES (
        v_event,
        (SELECT workflow FROM documents WHERE id = NEW.documents_id),
        NEW.changedate,
        NEW.actor,
        NEW.documents_id,
        NULL
    );
    RETURN NEW;
END;
$$;

CREATE TRIGGER ssvc_history_log_event
AFTER INSERT ON ssvc_history
FOR EACH ROW EXECUTE FUNCTION log_ssvc_history_to_events();


GRANT INSERT, DELETE, SELECT, UPDATE ON ssvc_history            TO {{ .User | sanitize }};
