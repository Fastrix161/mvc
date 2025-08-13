-- Menu
INSERT INTO Item (name, price, category, img) VALUES
    ('Paneer Tikka', '200','starters','images/placeHolder1.jpg'),
    ('Bread Omelete', '100','breakfast','images/placeHolder1.jpg'),
    ('Malai Kofte', '400','main course','images/placeHolder2.jpg'),
    ('Hot Coffee', '60','beverages','images/placeHolder2.jpg'),
    ('Chocolate Brownie', '120','dessert','images/placeHolder3.jpg'),
    ('Spring Rolls','200','starters','images/placeHolder1.jpg'),
    ('PBJ sandwitch','80','breakfast','images/placeHolder3.jpg'),
    ('Paneer Lababdaar','400','main course','images/placeHolder1.jpg'),
    ('Mojito','140','beverages','images/placeHolder2.jpg'),
    ('Ice Cream (Vanilla)','70','dessert','images/placeHolder3.jpg');
    
    -- User
    INSERT INTO User (name, email,role, mobile_number,password) VALUES
    ('Admin','admin@admin.com','admin','1010101010', '$2a$12$mTbQ9lM1YJK/DjT1/PBcoullmeBpzKC2ABlf22/OUPvgY5XHrAUVS'), --password: PasswordAdmin
    ('Chef','chef@gmail.com','chef','1231231230', '$2a$10$jgLHoHhS9WijLymbX.direU0MMMrzYmgxpQbSW6vF8vlwbGqGHrhy'); --password: PrimeChef
    -- SaltRound used: 10