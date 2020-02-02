SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for dynamicmaps
-- ----------------------------
DROP TABLE IF EXISTS `dynamicmaps`;
CREATE TABLE `dynamicmaps`  (
  `UUID` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `SerializedData` blob NULL,
  PRIMARY KEY (`UUID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for players
-- ----------------------------
DROP TABLE IF EXISTS `players`;
CREATE TABLE `players`  (
  `UUID` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `SerializedData` blob NULL,
  PRIMARY KEY (`UUID`) USING BTREE,
  UNIQUE INDEX `UUID_UNIQUE`(`UUID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles`  (
  `UUID` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `SummaryData` blob NULL,
  `OnlineData` blob NULL,
  PRIMARY KEY (`UUID`) USING BTREE,
  UNIQUE INDEX `UUID_UNIQUE`(`UUID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for staticmaps
-- ----------------------------
DROP TABLE IF EXISTS `staticmaps`;
CREATE TABLE `staticmaps`  (
  `UUID` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `SerializedData` blob NULL,
  `ZoneID` int(11) NULL DEFAULT NULL,
  `MapID` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`UUID`) USING BTREE,
  UNIQUE INDEX `UUID_UNIQUE`(`UUID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
