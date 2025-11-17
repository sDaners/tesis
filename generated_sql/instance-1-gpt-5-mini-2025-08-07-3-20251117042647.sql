CREATE TABLE Third_Party_Companies (
  company_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  company_type STRING(5) NOT NULL,
  company_name STRING(255),
  company_address STRING(255),
  other_company_details STRING(255)
) PRIMARY KEY (company_id);

CREATE TABLE Maintenance_Contracts (
  maintenance_contract_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  maintenance_contract_company_id STRING(36) NOT NULL REFERENCES Third_Party_Companies(company_id),
  contract_start_date TIMESTAMP,
  contract_end_date TIMESTAMP,
  other_contract_details STRING(255)
) PRIMARY KEY (maintenance_contract_id);

CREATE TABLE Parts (
  part_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  part_name STRING(255),
  chargeable_yn STRING(1),
  chargeable_amount NUMERIC,
  other_part_details STRING(255)
) PRIMARY KEY (part_id);

CREATE TABLE Skills (
  skill_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  skill_code STRING(20),
  skill_description STRING(255)
) PRIMARY KEY (skill_id);

CREATE TABLE Staff (
  staff_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  staff_name STRING(255),
  gender STRING(1),
  other_staff_details STRING(255)
) PRIMARY KEY (staff_id);

CREATE TABLE Assets (
  asset_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  maintenance_contract_id STRING(36) NOT NULL REFERENCES Maintenance_Contracts(maintenance_contract_id),
  supplier_company_id STRING(36) NOT NULL REFERENCES Third_Party_Companies(company_id),
  asset_details STRING(255),
  asset_make STRING(20),
  asset_model STRING(20),
  asset_acquired_date TIMESTAMP,
  asset_disposed_date TIMESTAMP,
  other_asset_details STRING(255)
) PRIMARY KEY (asset_id);

CREATE TABLE Asset_Parts (
  asset_id STRING(36) NOT NULL REFERENCES Assets(asset_id),
  part_id STRING(36) NOT NULL REFERENCES Parts(part_id)
) PRIMARY KEY (asset_id, part_id);

CREATE TABLE Maintenance_Engineers (
  engineer_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  company_id STRING(36) NOT NULL REFERENCES Third_Party_Companies(company_id),
  first_name STRING(50),
  last_name STRING(50),
  other_details STRING(255)
) PRIMARY KEY (engineer_id);

CREATE TABLE Engineer_Skills (
  engineer_id STRING(36) NOT NULL REFERENCES Maintenance_Engineers(engineer_id),
  skill_id STRING(36) NOT NULL REFERENCES Skills(skill_id)
) PRIMARY KEY (engineer_id, skill_id);

CREATE TABLE Fault_Log (
  fault_log_entry_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  asset_id STRING(36) NOT NULL REFERENCES Assets(asset_id),
  recorded_by_staff_id STRING(36) NOT NULL REFERENCES Staff(staff_id),
  fault_log_entry_datetime TIMESTAMP,
  fault_description STRING(255),
  other_fault_details STRING(255)
) PRIMARY KEY (fault_log_entry_id);

CREATE TABLE Engineer_Visits (
  engineer_visit_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  contact_staff_id STRING(36) REFERENCES Staff(staff_id),
  engineer_id STRING(36) NOT NULL REFERENCES Maintenance_Engineers(engineer_id),
  fault_log_entry_id STRING(36) NOT NULL REFERENCES Fault_Log(fault_log_entry_id),
  fault_status STRING(10) NOT NULL,
  visit_start_datetime TIMESTAMP,
  visit_end_datetime TIMESTAMP,
  other_visit_details STRING(255)
) PRIMARY KEY (engineer_visit_id);

CREATE TABLE Part_Faults (
  part_fault_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  part_id STRING(36) NOT NULL REFERENCES Parts(part_id),
  fault_short_name STRING(20),
  fault_description STRING(255),
  other_fault_details STRING(255)
) PRIMARY KEY (part_fault_id);

CREATE TABLE Fault_Log_Parts (
  fault_log_entry_id STRING(36) NOT NULL REFERENCES Fault_Log(fault_log_entry_id),
  part_fault_id STRING(36) NOT NULL REFERENCES Part_Faults(part_fault_id),
  fault_status STRING(10) NOT NULL
) PRIMARY KEY (fault_log_entry_id, part_fault_id);

CREATE TABLE Skills_Required_To_Fix (
  part_fault_id STRING(36) NOT NULL REFERENCES Part_Faults(part_fault_id),
  skill_id STRING(36) NOT NULL REFERENCES Skills(skill_id)
) PRIMARY KEY (part_fault_id, skill_id);