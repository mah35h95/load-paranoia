WITH
  field_ordered AS (
  SELECT
    *
  FROM
    `prod-2434-entdataingest-05104f.S4HANA.dd03l`
  ORDER BY
    tabname,
    position ),
  struct_it AS (
  SELECT
    CONCAT("{\nTableID: \"", LOWER(tabname),"\",\nColumns: []string{",CONCAT('"',STRING_AGG(LOWER(fieldname), '", "'),'"'),"},\n},") AS struct_value,
  FROM
    field_ordered
  WHERE
    LOWER(tabname) IN ( "bseg",
      "mldoc",
      "matdoc" )
    AND keyflag="X"
    AND LOWER(fieldname) != ".include"
  GROUP BY
    tabname )
SELECT
  string_agg(struct_value, "\n") as all_struct_value
FROM
  struct_it