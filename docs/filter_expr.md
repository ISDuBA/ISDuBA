# Filter expressions

ISDuBA supports filter expressions to narrow down on the advisories,
documents or events you are really looking for.

To do this you have to chain **conditions** in a [reverse polish notation](https://en.wikipedia.org/wiki/Reverse_Polish_notation) form. A **condition** is an boolean expression evaluating to `true` or `false`.

To select e.g. documents with an CVSSv3 greater or equal 5 you can write the following

```
$cvss_v3_score 5 float >=
```

Terms starting with `$` like `$cvss_v3_score` fetch data from the database. In this case the
value of the `cvss_v3_score` column of a document.

Each column from the database has a **data type**.  
The data type of the `$cvss_v3_score` is `float`.

To compare another value to this value the other value has to be of a compatible data type.
Writing simply `5` would be treated as string. To make it a float we have to **cast** it with
appending the required tyoe to it: `5 float`.  

Now `>=` as an **operator** can be applied to `$cvss_v3_score` as
both operands are type `float`. `>=` is a condition evaluating to `true` or `false`.

If you want another condition like the document should have a current release date after
the date 2023-12-31 you can extend this expression with

```
$cvss_v3_score 5 float >= $current_release_date 2023-12-31 time > and
```

`$current_release_date` results in a timestamp. `2023-12-31 time` too.
`>` checks for the order. `and` chains the CVSSv3 condition to the second one.
The `and` can be omitted in this case as all remainig conditions which are
not explicity connected are `and`ed together.


See the [Examples](#section_examples) section for more examples.  
See the [Columns](#section_columns) section for which data fields are available.  
See the [Operators](#section_operators) section for the available operatores.  
See the [Data types](#section_datatypes) section for the available data types.  

## <a name="section_examples"></a> Examples

**TBD**

## <a name="section_columns"></a> Columns

| Column                 | Data type | Document | Advisory | Event | Description |
| ---------------------- | --------- | -------- | -------- | ----- | ----------- |
| `id`                   | integer   |         |         |      | Database ID of a document |
| `latest`               | bool      |         |         |      | Latest document of an advisory |
| `tracking_id`          | string    |         |         |      | `/document/tracking/id` |
| `version`              | string    |         |         |      | `/document/tracking/version` |
| `publisher`            | string    |         |         |      | `/document/publisher/name` |
| `current_release_date` | time      |         |         |      | `/document/tracking/current_release_date` |
| `initial_release_date` | time      |         |         |      | `/document/tracking/initial_release_date` |
| `rev_history_length`   | int       |         |         |      | Length of the revision history |
| `title`                | string    |         |         |      | `/document/title` |
| `tlp`                  | string    |         |         |      | `/document/distribution/tlp/label` |
| `ssvc`                 | string    |         |         |      | SSVC score of this document |
| `cvss_v2_score`        | float     |         |         |      | `max(/document/vulnerabilities[*]/scores[*]/cvss_v2/baseScore)` |
| `cvss_v3_score`        | float     |         |         |      | `max(/document/vulnerabilities[*]/scores[*]/cvss_v3_scorecore)` |
| `critical`             | float     |         |         |      | `coalesce(cvss_v3_score, cvss_v2_score)` |
| `comments`             | int       |         |         |      | Number of comments of document/advisory |
| `state`                | workflow  |         |         |      | State of advisory |
| `recent`               | time      |         |         |      | Timestamp of recent event of advisory |
| `versions`             | int       |         |         |      | Number of documents per advisory |
| `event`                | events    |         |         |      | Type of event |
| `event_state`          | workflow  |         |         |      | State of advisory associated with event |
| `time`                 | time      |         |         |      | Timestamp of the event |
| `actor`                | string    |         |         |      | User who triggered the event |
| `comments_id`          | int       |         |         |      | If event was comment related, ID of the affected comment |

## <a name="section_examples"></a> Examples

**TBD**


**TBD**

## <a name="section_operators"></a> Operators

**TBD**

## <a name="section_datatypes"></a>Data types

- `float`: Floating point numbers
- `integer`: Integer numbers
- `bool`: Boolean values
- `string`: String/Text values.
- `time`: Timestamps
- `duration`: Length of time intervals.
- `workflow`: States of workflow.
- `events`: States of events.

