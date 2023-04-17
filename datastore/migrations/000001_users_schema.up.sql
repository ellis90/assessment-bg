CREATE TABLE "users" (
                                "id" bigserial PRIMARY KEY,
                                "user_name" varchar(50) UNIQUE,
                                "first_name" varchar(255),
                                "last_name" varchar(255),
                                "email" varchar(255) UNIQUE,
                                "department" varchar(255),
                                "user_status" varchar(1)
);