package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upInit, downInit)
}

func upInit(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create schema if not exists service;

		create table service.users (
			id uuid not null primary key,
			first_name varchar(255) not null,
			second_name varchar(255) not null,
			email varchar(255) not null unique
		);

		create table service.products (
			id uuid not null primary key,
			name varchar(255) not null,
			price bigint default 0 not null,
			duration_days int not null
		);

		create type subscription_status as enum ('active', 'paused', 'canceled');
		create table service.subscriptions  (
			id uuid not null primary key,
			user_id uuid not null,
    		product_id uuid not null,
			start_date timestamp not null,
			end_date timestamp not null,
			duration_days int not null,
			price bigint default 0 not null,
			tax bigint default 0 not null,
			total_price bigint default 0 not null,
			total_price_with_voucher bigint default 0 not null,
			status subscription_status not null,
			trial_start_date timestamp,
			trial_end_date timestamp,
			canceled_date timestamp,
			paused_date timestamp,
			foreign key (user_id) references service.users(id),
    		foreign key (product_id) references service.products(id)
		);

		create type voucher_status as enum ('percentage', 'fixed');
		create table service.vouchers (
			id uuid not null primary key,
			code varchar(255) not null unique,
			discount_type voucher_status not null,
			discount_value bigint default 0 not null,
			valid_from timestamp not null,
			valid_until timestamp not null
		);

		create table service.product_vouchers (
			product_id uuid not null,
			voucher_id uuid not null,
			primary key (product_id, voucher_id),
			foreign key (product_id) references service.products(id),
			foreign key (voucher_id) references service.vouchers(id)
		);

		create table service.subscriptions_vouchers (
			subscription_id uuid not null,
			voucher_id uuid not null,
			primary key (subscription_id, voucher_id),
			foreign key (subscription_id) references service.subscriptions(id),
			foreign key (voucher_id) references service.vouchers(id)
		);

		insert into service.users (id, first_name, second_name, email) values
			('b5d2f6ec-5eac-4e62-8ac0-3c45e1b9f3b5', 'John', 'Doe', 'john.doe@example.com'),
			('4e8f3d24-5e1e-4a96-9a2d-9d65b4bc07da', 'Jane', 'Smith', 'jane.smith@example.com'),
			('1f3f4c0c-088d-4784-803a-2b7e4679c4a8', 'Michael', 'Johnson', 'michael.johnson@example.com')
		;
		
		insert into service.products (id, name, price, duration_days) values
			('a72d8c5c-cb57-42d2-b3b2-13e9ed06403b', 'Basic Plan', 1000, 30),
			('9b1e0b0b-3c34-4cfa-8f63-5d12b3feff34', 'Standart Plan', 1500, 60),
			('ab97234d-6b4a-4a70-823e-68b7a80ef6d4', 'Premium Plan', 2500, 90),
			('29fdcb93-b52f-48a9-9e7e-b3e60d63d8a3', 'Enterprise Plan', 35000, 365)
		;

		insert into service.subscriptions (id, user_id, product_id, start_date, end_date, duration_days, price, tax, total_price, total_price_with_voucher, status, trial_start_date, trial_end_date) values
			('0e1f7de1-e9c1-4bb9-bb63-69bba837b706', 'b5d2f6ec-5eac-4e62-8ac0-3c45e1b9f3b5', 'a72d8c5c-cb57-42d2-b3b2-13e9ed06403b', '2025-03-01 10:00:00', '2025-03-31 10:00:00', 30, 5000, 500, 5500, 5500, 'active', '2025-03-01 10:00:00', '2025-03-01 10:00:00'),
			('d734e51b-8972-4d3f-b7f5-e7bfb24e9b2a', '4e8f3d24-5e1e-4a96-9a2d-9d65b4bc07da', 'ab97234d-6b4a-4a70-823e-68b7a80ef6d4', '2025-03-10 12:00:00', '2025-06-10 12:00:00', 90, 15000, 1500, 16500, 16500, 'active', '2025-03-10 12:00:00', '2025-03-10 12:00:00'),
			('6f8a0d98-8df0-44f9-b401-cde8ccbf63a0', '1f3f4c0c-088d-4784-803a-2b7e4679c4a8', '29fdcb93-b52f-48a9-9e7e-b3e60d63d8a3', '2025-01-01 09:00:00', '2025-12-31 09:00:00', 365, 35000, 3500, 38500, 38500, 'active', '2025-01-01 09:00:00', '2025-01-01 09:00:00')
		;
	
		insert into service.vouchers (id, code, discount_type, discount_value, valid_from, valid_until) values
			('c1b0e227-17bb-4335-a1c2-e61c6e635b84', 'DISCOUNT10', 'percentage', 10, '2025-03-01 00:00:00', '2025-03-31 23:59:59'),
			('50f12d5d-0587-4e88-b7e3-bde71737d7ef', 'FIXED5', 'fixed', 500, '2025-03-01 00:00:00', '2025-03-31 23:59:59'),
			('df0a1fd1-f560-49b9-b287-b4c699710c1f', 'SUMMER25', 'percentage', 25, '2025-06-01 00:00:00', '2025-06-30 23:59:59')
		;
		
		insert into service.product_vouchers (product_id, voucher_id) values
			('a72d8c5c-cb57-42d2-b3b2-13e9ed06403b', 'c1b0e227-17bb-4335-a1c2-e61c6e635b84'),
			('ab97234d-6b4a-4a70-823e-68b7a80ef6d4', '50f12d5d-0587-4e88-b7e3-bde71737d7ef'),
			('29fdcb93-b52f-48a9-9e7e-b3e60d63d8a3', 'df0a1fd1-f560-49b9-b287-b4c699710c1f')
		;

		insert into service.subscriptions_vouchers (subscription_id, voucher_id) values
			('0e1f7de1-e9c1-4bb9-bb63-69bba837b706', 'c1b0e227-17bb-4335-a1c2-e61c6e635b84'),
			('d734e51b-8972-4d3f-b7f5-e7bfb24e9b2a', '50f12d5d-0587-4e88-b7e3-bde71737d7ef'),
			('6f8a0d98-8df0-44f9-b401-cde8ccbf63a0', 'df0a1fd1-f560-49b9-b287-b4c699710c1f')
		;
	`)
	if err != nil {
		return err
	}

	return nil
}

func downInit(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop table if exists service.users;
		drop table if exists service.subscriptions;
		drop table if exists service.vouchers;
		drop table if exists service.products;
		drop schema if exists service cascade;
	`)
	if err != nil {
		return err
	}

	return nil
}
