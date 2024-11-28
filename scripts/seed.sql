-- Generate sample customers
INSERT INTO
    Customer (
        CustomerID,
        ClientCustomerID,
        InsertDate
    )
SELECT n, n * 10, DATE_SUB(
        CURRENT_TIMESTAMP, INTERVAL FLOOR(RAND() * 365) DAY
    )
FROM (
        SELECT ROW_NUMBER() OVER () as n
        FROM information_schema.columns
        LIMIT 1000
    ) nums;

-- Generate customer emails
INSERT INTO
    CustomerData (
        CustomerChannelID,
        CustomerID,
        ChannelTypeID,
        ChannelValue,
        InsertDate
    )
SELECT
    CustomerID,
    CustomerID,
    1, -- Email type
    CONCAT(
        'customer_',
        CustomerID,
        '@example.com'
    ),
    InsertDate
FROM Customer;

-- Generate sample content
INSERT INTO
    Content (
        ContentID,
        ClientContentID,
        InsertDate
    )
SELECT n, n * 100, CURRENT_TIMESTAMP
FROM (
        SELECT ROW_NUMBER() OVER () as n
        FROM information_schema.columns
        LIMIT 100
    ) nums;

-- Generate content prices
INSERT INTO
    ContentPrice (
        ContentPriceID,
        ContentID,
        Price,
        Currency,
        InsertDate
    )
SELECT ContentID, ContentID, ROUND(10 + RAND() * 990, 2), -- Random price between 10 and 1000
    'USD', CURRENT_TIMESTAMP
FROM Content;

-- Generate customer events
INSERT INTO
    CustomerEvent (
        EventID,
        ClientEventID,
        InsertDate
    )
SELECT
    ROW_NUMBER() OVER () as EventID,
    ROW_NUMBER() OVER () * 1000 as ClientEventID,
    CURRENT_TIMESTAMP
FROM (
        SELECT a.n
        FROM (
                SELECT ROW_NUMBER() OVER () as n
                FROM information_schema.columns
                LIMIT 100
            ) a, (
                SELECT ROW_NUMBER() OVER () as n
                FROM information_schema.columns
                LIMIT 100
            ) b
        LIMIT 5000
    ) nums;

-- Generate purchase events
INSERT INTO
    CustomerEventData (
        EventDataID,
        EventID,
        ContentID,
        CustomerID,
        EventTypeID,
        EventDate,
        Quantity,
        InsertDate
    )
SELECT
    ROW_NUMBER() OVER () as EventDataID,
    e.EventID,
    FLOOR(1 + RAND() * 100) as ContentID, -- Random content
    FLOOR(1 + RAND() * 1000) as CustomerID, -- Random customer
    6, -- Purchase event type
    DATE_ADD(
        '2020-04-01',
        INTERVAL FLOOR(RAND() * 365) DAY
    ) as EventDate,
    FLOOR(1 + RAND() * 5) as Quantity, -- Random quantity between 1 and 5
    CURRENT_TIMESTAMP
FROM CustomerEvent e;