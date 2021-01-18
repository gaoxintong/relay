INSERT INTO `iot-hub`.`product` (
    `product_name`,
    `product_description`,
    `product_key`,
    `product_config`,
    `product_image`
  )
VALUES (
    'MI8继电器',
    '物业设备 MI8继电器',
    'abdf6b26a399494869c5db5476d1d617fdb5f7d6579fd093ccf78c77ea61e70f',
    '{\n  \"objects\": [\n    {\n      \"no\": 1,\n      \"label\": \"switch1\",\n      \"part\": 1,\n      \"status\": [\n        {\n          \"value_type\": 8,\n          \"name\": \"status\"\n        }\n      ]\n    },\n    {\n      \"no\": 2,\n      \"label\": \"switch2\",\n      \"part\": 1,\n      \"status\": [\n        {\n          \"value_type\": 8,\n          \"name\": \"status\"\n        }\n      ]\n    },\n    {\n      \"no\": 3,\n      \"label\": \"switch3\",\n      \"part\": 1,\n      \"status\": [\n        {\n          \"value_type\": 8,\n          \"name\": \"status\"\n        }\n      ]\n    },\n    {\n      \"no\": 4,\n      \"label\": \"switch4\",\n      \"part\": 1,\n      \"status\": [\n        {\n          \"value_type\": 8,\n          \"name\": \"status\"\n        }\n      ]\n    },\n    {\n      \"no\": 5,\n      \"label\": \"switch5\",\n      \"part\": 1,\n      \"status\": [\n        {\n          \"value_type\": 8,\n          \"name\": \"status\"\n        }\n      ]\n    },\n    {\n      \"no\": 6,\n      \"label\": \"switch6\",\n      \"part\": 1,\n      \"status\": [\n        {\n          \"value_type\": 8,\n          \"name\": \"status\"\n        }\n      ]\n    },\n    {\n      \"no\": 7,\n      \"label\": \"switch7\",\n      \"part\": 1,\n      \"status\": [\n        {\n          \"value_type\": 8,\n          \"name\": \"status\"\n        }\n      ]\n    },\n    {\n      \"no\": 8,\n      \"label\": \"switch8\",\n      \"part\": 1,\n      \"status\": [\n        {\n          \"value_type\": 8,\n          \"name\": \"status\"\n        }\n      ]\n    }\n  ],\n  \"commands\": [\n    {\n      \"no\": 1,\n      \"name\": \"setState\",\n      \"part\": 1,\n      \"priority\": 0,\n      \"params\": [\n        {\n          \"value_type\": 7,\n          \"name\": \"no\"\n        },\n        {\n          \"value_type\": 7,\n          \"name\": \"state\"\n        }\n      ]\n    }\n  ],\n  \"events\": []\n}',
    ''
  );