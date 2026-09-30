// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package database

import (
	"context"

	"github.com/ISDuBA/ISDuBA/pkg/models"
)

func init() {
	RegisterMigrationCode("update_documents_index", func(ctx context.Context, t Transaction) error {
		tx, err := t.Transaction(ctx)
		if err != nil {
			return err
		}
		return models.UpdateDocumentsIndex(ctx, tx)
	})
}
