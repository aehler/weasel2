create table emt.market_history (
	"typeID" bigint not null references evesde."invTypes" ("typeID"),
	"regionID" bigint not null references evesde."mapRegions" ("regionID"),
	dt date not null,
	order_count bigint not null,
	volume bigint not null,
	highest numeric not null,
	average numeric not null,
	lowest numeric not null
);

create index idx_market_history_type_id on emt.market_history ("typeID");
create index idx_market_history_region_id on emt.market_history ("regionID");
create index idx_market_history_dt on emt.market_history ("dt");

create table emt.market_history_updates (
  "typeID" bigint not null REFERENCES evesde."invTypes" ("typeID"),
  "regionID" bigint not null references evesde."mapRegions" ("regionID"),
  date_update timestamp default CURRENT_TIMESTAMP
);

create or replace function emt.market_history_upsert (_type_id bigint, _region_id bigint) returns void as $$
  begin
      WITH upsert AS (
            UPDATE emt.market_history_updates
                SET
                    date_update    = current_timestamp
            WHERE "typeID" = _type_id and "regionID" = _region_id
            RETURNING *
        )
        INSERT INTO emt.market_history_updates ("typeID", "regionID")
            SELECT _type_id, _region_id WHERE NOT EXISTS (SELECT * FROM upsert);
  end;
$$ language plpgsql;
