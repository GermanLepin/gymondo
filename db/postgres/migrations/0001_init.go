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
			duration_days int not null,
			price decimal(15,2) default 0 not null,
			tax decimal(15,2) default 0 not null,
			total_price decimal(15,2) default 0 not null
		);
		
		create type subscription_status as enum ('active', 'paused', 'canceled');
		create table service.subscriptions (
			id uuid not null primary key,
			user_id uuid references service.users(id) on delete cascade,
			product_id uuid references service.products(id) on delete cascade,
			start_date timestamp not null,
			end_date timestamp not null,
			duration_days int not null,
			price decimal(15,2) default 0 not null,
			tax decimal(15,2) default 0 not null,
			total_price decimal(15,2) default 0 not null,
			status subscription_status not null,
			trial_start_date timestamp,
			trial_end_date timestamp,
			canceled_date timestamp,
			paused_date timestamp,
			unpaused_date timestamp
		);
		
		create type voucher_status as enum ('percentage', 'fixed');
		create table service.vouchers (
			id uuid not null primary key,
			code varchar(255) not null unique,
			discount_type voucher_status not null,
			discount_value decimal(5,2) default 0 not null
		);
		
		insert into service.users (id, first_name, second_name, email) values
			('b5d2f6ec-5eac-4e62-8ac0-3c45e1b9f3b5', 'john', 'doe', 'john.doe@example.com'),
			('4e8f3d24-5e1e-4a96-9a2d-9d65b4bc07da', 'jane', 'smith', 'jane.smith@example.com'),
			('1f3f4c0c-088d-4784-803a-2b7e4679c4a8', 'michael', 'johnson', 'michael.johnson@example.com')
		;
		
		insert into service.products (id, name, duration_days, price, tax, total_price) values
			('a72d8c5c-cb57-42d2-b3b2-13e9ed06403b', 'basic plan',  30, 10.00, 1.00, 11.00),
			('9b1e0b0b-3c34-4cfa-8f63-5d12b3feff34', 'standard plan', 60, 15.00, 1.50, 16.50),
			('ab97234d-6b4a-4a70-823e-68b7a80ef6d4', 'premium plan', 90, 25.00, 2.50, 27.50),
			('29fdcb93-b52f-48a9-9e7e-b3e60d63d8a3', 'enterprise plan', 365, 80.00, 8.00, 88.00)
		;
		
		insert into service.vouchers (id, code, discount_type, discount_value) values
			('c1b0e227-17bb-4335-a1c2-e61c6e635b84', 'discount10', 'percentage', 0.10),
			('50f12d5d-0587-4e88-b7e3-bde71737d7ef', 'fixed5', 'fixed', 5.00),
			('df0a1fd1-f560-49b9-b287-b4c699710c1f', 'summer25', 'percentage', 0.25)
		;
	`)
	if err != nil {
		return err
	}

	return nil
}

func downInit(tx *sql.Tx) error {
	return nil
}
