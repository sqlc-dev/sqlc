CREATE TABLE group_calc_totals (
  npn Text,
  group_id Text,
  PRIMARY KEY (group_id)
);

CREATE TABLE producer_group_attribute (
  npn_external_map_id Text,
  group_id Text,
  PRIMARY KEY (group_id)
);

CREATE TABLE npn_external_map (
  id Text,
  npn Text,
  PRIMARY KEY (id)
);
