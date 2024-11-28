CREATE TABLE ChannelType (
    ChannelTypeID SMALLINT PRIMARY KEY,
    Name VARCHAR(30)
);

CREATE TABLE CustomerData (
    CustomerChannelID BIGINT PRIMARY KEY,
    CustomerID BIGINT,
    ChannelTypeID SMALLINT,
    ChannelValue VARCHAR(600),
    InsertDate TIMESTAMP,
    FOREIGN KEY (ChannelTypeID) REFERENCES ChannelType (ChannelTypeID)
);

CREATE TABLE Customer (
    CustomerID BIGINT PRIMARY KEY,
    ClientCustomerID BIGINT,
    InsertDate TIMESTAMP
);

CREATE TABLE EventType (
    EventTypeID SMALLINT PRIMARY KEY,
    Name VARCHAR(30)
);

CREATE TABLE CustomerEvent (
    EventID BIGINT PRIMARY KEY,
    ClientEventID BIGINT,
    InsertDate TIMESTAMP
);

CREATE TABLE CustomerEventData (
    EventDataID BIGINT PRIMARY KEY,
    EventID BIGINT,
    ContentID INT,
    CustomerID BIGINT,
    EventTypeID SMALLINT,
    EventDate TIMESTAMP,
    Quantity SMALLINT,
    InsertDate TIMESTAMP,
    FOREIGN KEY (EventID) REFERENCES CustomerEvent (EventID),
    FOREIGN KEY (CustomerID) REFERENCES Customer (CustomerID),
    FOREIGN KEY (EventTypeID) REFERENCES EventType (EventTypeID)
);

CREATE TABLE Content (
    ContentID INT PRIMARY KEY,
    ClientContentID BIGINT,
    InsertDate TIMESTAMP
);

CREATE TABLE ContentPrice (
    ContentPriceID BIGINT PRIMARY KEY,
    ContentID INT,
    Price DECIMAL(8, 2),
    Currency CHAR(3),
    InsertDate TIMESTAMP,
    FOREIGN KEY (ContentID) REFERENCES Content (ContentID)
);

-- Insert channel types
INSERT INTO
    ChannelType
VALUES (1, 'Email'),
    (2, 'PhoneNumber'),
    (3, 'Postal'),
    (4, 'MobileID'),
    (5, 'Cookie');

-- Insert event types
INSERT INTO
    EventType
VALUES (1, 'sent'),
    (2, 'view'),
    (3, 'click'),
    (4, 'visit'),
    (5, 'cart'),
    (6, 'purchase');