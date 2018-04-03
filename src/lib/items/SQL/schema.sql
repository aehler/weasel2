drop materialized view emt.production_efficiency;

create materialized view emt.production_efficiency as (

  with tt as (
  select tp."typeName", prod."typeID", prod."productTypeID" as product_type_id, g."groupName", array_agg(m."materialTypeID") as mat_ids, tpp."typeName" as product, prod.quantity,
  avg(mpc.avg_price * prod.quantity) as prod_sell_price, sum(mpc2.avg_price * m.quantity) as material_cost, g_prod."groupName" as prod_group_name,
  mg."marketGroupName" is not null as bpo_purchasable,
  mtg."metaGroupName", mtg."metaGroupID",
  (select coalesce("valueInt", "valueFloat") from evesde."dgmTypeAttributes" as mt where mt."typeID" = tpp."typeID" and mt."attributeID" = 422) as tech_lvl,
  (select coalesce("valueInt", "valueFloat") from evesde."dgmTypeAttributes" as mt where mt."typeID" = tpp."typeID" and mt."attributeID" = 633) as meta_lvl,
  (select exists (select 1 from evesde."industryActivityProducts" where "productTypeID" = bpo."typeID")) as obtained_by_research,
  array_agg(mpc2.avg_price) as avg_prices,
  json_agg(json_build_object('material_type_id', m."materialTypeID", 'material', tpm."typeName", 'avg_price', coalesce(mpc2.avg_price, 0), 'quantity', m.quantity)) as mats_json
  from evesde."industryBlueprints" as bpo
  left join evesde."invTypes" as tp using("typeID")
  left join evesde."invGroups" as g on g."groupID" = tp."groupID"
  left join evesde."industryActivityMaterials" as m on m."typeID" = bpo."typeID" and m."activityID" = 1
  left join evesde."invTypes" as tpm on m."materialTypeID" = tpm."typeID"
  left join evesde."industryActivityProducts" as prod on bpo."typeID" = prod."typeID"
  left join evesde."invTypes" as tpp on prod."productTypeID" = tpp."typeID"
  left join evesde."invGroups" as g_prod on tpp."groupID" = g_prod."groupID"
  left join evesde."invMarketGroups" as mg on mg."marketGroupID" = tp."marketGroupID"
  left join emt.market_avg as mpc on mpc.type_id = tpp."typeID"
  left join emt.market_avg as mpc2 on mpc2.type_id = tpm."typeID"
  left join evesde."invMetaGroups" as mtg on (select coalesce("valueInt", "valueFloat") from evesde."dgmTypeAttributes" as mt where mt."typeID" = prod."productTypeID" and mt."attributeID" = 1692) = mtg."metaGroupID"
  where m."activityID" = 1 and tpp.published = true
  and mpc.avg_price is not null
  group by tp."typeName", g."groupName", prod."productTypeID", tpp."typeName", prod.quantity, prod."typeID", bpo."typeID",
  g_prod."groupName", bpo_purchasable, mtg."metaGroupName", tech_lvl, meta_lvl, mtg."metaGroupID")

  select tt.prod_sell_price - material_cost as net_profit,  case tt.material_cost when 0 then 0 else tt.prod_sell_price / material_cost end as ROS,
  to_tsvector(tt."product") || to_tsvector(tt.prod_group_name) as norm_vc, tt.* from tt
  where tt.product_type_id != any(tt.mat_ids)

)
with no data;

CREATE INDEX idx_emt_production_efficiency_fts ON emt.production_efficiency
  USING gin(norm_vc);

refresh materialized view emt.production_efficiency;

create materialized view emt.market_averages as (

  with tt as

) with no data;



create materialized view emt.avg_grouped as (
select count(*) as days, "regionID", "typeID", round(avg(order_count)) as order_count, round(avg(volume)) as volume,
max(dt) as latest_dt, current_timestamp + '-90 days' as back_dt,
DATE_PART('day', max(dt)::timestamp without time zone - (current_timestamp + '-90 days')::timestamp without time zone) as days_back
from emt.market_history where dt > current_timestamp + '-90 days'
group by "regionID", "typeID"
order by "typeID"
) with no data;

refresh materialized view emt.avg_grouped;