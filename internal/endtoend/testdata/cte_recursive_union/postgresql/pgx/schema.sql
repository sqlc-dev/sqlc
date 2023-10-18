CREATE TABLE case_intent_version
(
    version_id SERIAL NOT NULL PRIMARY KEY,
    reviewer TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE TABLE case_intent
(
    id SERIAL NOT NULL PRIMARY KEY,
    case_intent_string TEXT NOT NULL,
    description TEXT NOT NULL,
    author TEXT NOT NULL
);
CREATE TABLE case_intent_parent_join
(
    case_intent_id BIGINT NOT NULL,
    case_intent_parent_id BIGINT NOT NULL,
    constraint fk_case_intent_id foreign key (case_intent_id) references case_intent(id),
    constraint fk_case_intent_parent_id foreign key (case_intent_parent_id) references case_intent(id)
);
CREATE TABLE case_intent_version_join
(
    case_intent_id BIGINT NOT NULL,
    case_intent_version_id INT NOT NULL,
    constraint fk_case_intent_id foreign key (case_intent_id) references case_intent(id),
    constraint fk_case_intent_version_id foreign key (case_intent_version_id) references case_intent_version(version_id)
);
