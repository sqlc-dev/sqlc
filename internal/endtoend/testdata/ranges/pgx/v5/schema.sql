CREATE TABLE test_table
(
    v_daterange_null      daterange,
    v_datemultirange_null datemultirange,
    v_tsrange_null        tsrange,
    v_tsmultirange_null   tsmultirange,
    v_tstzrange_null      tstzrange,
    v_tstzmultirange_null tstzmultirange,
    v_numrange_null       numrange,
    v_nummultirange_null  nummultirange,
    v_int4range_null      int4range,
    v_int4multirange_null int4multirange,
    v_int8range_null      int8range,
    v_int8multirange_null int8multirange,
    v_daterange           daterange      not null,
    v_datemultirange      datemultirange not null,
    v_tsrange             tsrange        not null,
    v_tsmultirange        tsmultirange   not null,
    v_tstzrange           tstzrange      not null,
    v_tstzmultirange      tstzmultirange not null,
    v_numrange            numrange       not null,
    v_nummultirange       nummultirange  not null,
    v_int4range           int4range      not null,
    v_int4multirange      int4multirange not null,
    v_int8range           int8range      not null,
    v_int8multirange      int8multirange not null
);

