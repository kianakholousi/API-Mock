
CREATE TABLE IF NOT EXISTS cities (
    id int PRIMARY KEY AUTO_INCREMENT,
    name varchar(255) UNIQUE NOT NULL ,
    created_at datetime DEFAULT NOW(),
    updated_at datetime DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS airplanes (
    id int PRIMARY KEY AUTO_INCREMENT,
    name varchar(255) UNIQUE NOT NULL ,
    created_at datetime DEFAULT NOW(),
    updated_at datetime DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS canceling_situations (
    id int PRIMARY KEY AUTO_INCREMENT,
    description varchar(255) NOT NULL ,
    data varchar(255) NOT NULL ,
    created_at datetime DEFAULT NOW(),
    updated_at datetime DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS flights (
    id int PRIMARY KEY AUTO_INCREMENT,
    dep_city_id int NOT NULL ,
    arr_city_id int NOT NULL ,
    dep_time datetime NOT NULL ,
    arr_time datetime NOT NULL ,
    airplane_id int NOT NULL ,
    airline varchar(255) NOT NULL ,
    price int NOT NULL ,
    cxl_sit_id int NOT NULL ,
    remaining_seats int NOT NULL,
    flight_class varchar(255) NOT NULL ,
    baggage_allowance int NOT NULL ,
    meal_service varchar(255) NOT NULL,
    gate_number int NOT NULL ,
    created_at datetime DEFAULT NOW(),
    updated_at datetime DEFAULT NOW() ON UPDATE NOW(),

    FOREIGN KEY (dep_city_id) REFERENCES cities(id),
    FOREIGN KEY (arr_city_id) REFERENCES cities(id),
    FOREIGN KEY (airplane_id) REFERENCES airplanes(id),
    FOREIGN KEY (cxl_sit_id) REFERENCES canceling_situations(id)
);
