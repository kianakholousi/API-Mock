
CREATE TABLE IF NOT EXISTS city (
    id int PRIMARY KEY AUTO_INCREMENT,
    name varchar(255) UNIQUE NOT NULL ,
    created_at datetime DEFAULT NOW(),
    updated_at datetime DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS airplane (
    id int PRIMARY KEY AUTO_INCREMENT,
    name varchar(255) UNIQUE NOT NULL ,
    created_at datetime DEFAULT NOW(),
    updated_at datetime DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS canceling_situation (
    id int PRIMARY KEY AUTO_INCREMENT,
    description varchar(255) NOT NULL ,
    data varchar(255) NOT NULL ,
    created_at datetime DEFAULT NOW(),
    updated_at datetime DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS flight (
    id int PRIMARY KEY AUTO_INCREMENT,
    dep_city_id int NOT NULL ,
    arr_city_id int NOT NULL ,
    dep_time datetime NOT NULL ,
    arr_time datetime NOT NULL ,
    airplane_id int NOT NULL ,
    airline varchar(255) NOT NULL ,
    price int NOT NULL ,
    cxl_sit_id int NOT NULL ,
    left_seat int NOT NULL,
    created_at datetime DEFAULT NOW(),
    updated_at datetime DEFAULT NOW() ON UPDATE NOW(),

    FOREIGN KEY (dep_city_id) REFERENCES city(id),
    FOREIGN KEY (arr_city_id) REFERENCES city(id),
    FOREIGN KEY (airplane_id) REFERENCES airplane(id),
    FOREIGN KEY (cxl_sit_id) REFERENCES canceling_situation(id)
);