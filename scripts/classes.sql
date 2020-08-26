create table classes(
   id INT NOT NULL AUTO_INCREMENT,
   class_name VARCHAR(100) NOT NULL,
   start_date DATE NOT NULL,
   end_date DATE NOT NULL,
   capacity INT,
   PRIMARY KEY ( id )
);