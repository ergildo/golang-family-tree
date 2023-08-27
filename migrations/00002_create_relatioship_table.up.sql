CREATE TABLE relationship
(
    parentId integer NOT NULL,
    childId integer NOT NULL,
    CONSTRAINT pk_relationship PRIMARY KEY (parentId, childId),
    CONSTRAINT fk_relationship_child FOREIGN KEY (childId)
        REFERENCES person (id) MATCH SIMPLE
        ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT fk_relationship_parent FOREIGN KEY (parentId)
        REFERENCES person (id) MATCH SIMPLE
        ON UPDATE CASCADE ON DELETE CASCADE
);