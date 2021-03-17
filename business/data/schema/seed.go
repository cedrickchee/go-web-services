package schema

import "github.com/jmoiron/sqlx"

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// seeds is a string constant containing all of the queries needed to get the
// db seeded to a useful state for development.
//
// Note that database servers besides PostgreSQL may not support running
// multiple queries as part of the same execution so this single large constant
// may need to be broken up.
const seeds = `
-- Create admin and regular User with password "sup3rS3cr3tGolang"
INSERT INTO users (user_id, name, email, roles, password_hash, date_created, date_updated) VALUES
	('1110d50a-7cea-4279-9f7d-b79ed659d31b', 'Admin', 'admin@example.com', '{ADMIN,USER}', '$2a$10$U9UWrPDsiX0oItb/Z8chJuJRsAPPn4SakhdkCqA4LTLMlfZaAgmIe', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('dadd74d1-55f6-4d3b-b203-b2898464ac91', 'User', 'user@example.com', '{USER}', '$2a$10$nDfc4WMyXoO6Pul19UYMi.AuxLVuDtEmpKPnolRm9vGnvN5jbE67e', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;

INSERT INTO products (product_id, user_id, name, cost, quantity, date_created, date_updated) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'dadd74d1-55f6-4d3b-b203-b2898464ac91', 'Bag', 50, 42, '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'dadd74d1-55f6-4d3b-b203-b2898464ac91', 'Shoes', 75, 120, '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;

INSERT INTO sales (sale_id, product_id, quantity, paid, date_created) VALUES
	('98b6d4b8-f04b-4c79-8c2e-a0aef46854b7', 'a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 2, 100, '2019-01-01 00:00:03.000001+00'),
	('85f6fb09-eb05-4874-ae39-82d1a30fe0d7', 'a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 5, 250, '2019-01-01 00:00:04.000001+00'),
	('a235be9e-ab5d-44e6-a987-fa1c749264c7', '72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 3, 225, '2019-01-01 00:00:05.000001+00')
	ON CONFLICT DO NOTHING;
`

// DeleteAll runs the set of Drop-table queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func DeleteAll(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(deleteAll); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

// deleteAll is used to clean the database between tests.
const deleteAll = `
DELETE FROM sales;
DELETE FROM products;
DELETE FROM users;
`
