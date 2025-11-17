CREATE TABLE Third_Party_Companies (
  company_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  company_type STRING(5) NOT NULL,
  company_name STRING(255),
  company_address STRING(255),
  other_company_details STRING(255)
) PRIMARY KEY (company_id);

CREATE TABLE Parts (
  part_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  part_name STRING(255),
  chargeable_yn STRING(1),
  chargeable_amount STRING(20),
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

CREATE TABLE Maintenance_Contracts (
  maintenance_contract_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  maintenance_contract_company_id STRING(36) NOT NULL,
  contract_start_date TIMESTAMP,
  contract_end_date TIMESTAMP,
  other_contract_details STRING(255),
  CONSTRAINT fk_maintenance_contract_company FOREIGN KEY (maintenance_contract_company_id)
    REFERENCES Third_Party_Companies(company_id)
) PRIMARY KEY (maintenance_contract_id);

CREATE TABLE Assets (
  asset_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  maintenance_contract_id STRING(36) NOT NULL,
  supplier_company_id STRING(36) NOT NULL,
  asset_details STRING(255),
  asset_make STRING(20),
  asset_model STRING(20),
  asset_acquired_date TIMESTAMP,
  asset_disposed_date TIMESTAMP,
  other_asset_details STRING(255),
  CONSTRAINT fk_assets_maintenance_contract FOREIGN KEY (maintenance_contract_id)
    REFERENCES Maintenance_Contracts(maintenance_contract_id),
  CONSTRAINT fk_assets_supplier_company FOREIGN KEY (supplier_company_id)
    REFERENCES Third_Party_Companies(company_id)
) PRIMARY KEY (asset_id);

CREATE TABLE Asset_Parts (
  asset_id STRING(36) NOT NULL,
  part_id STRING(36) NOT NULL,
  CONSTRAINT fk_asset_parts_asset FOREIGN KEY (asset_id) REFERENCES Assets(asset_id),
  CONSTRAINT fk_asset_parts_part FOREIGN KEY (part_id) REFERENCES Parts(part_id)
) PRIMARY KEY (asset_id, part_id);

CREATE TABLE Maintenance_Engineers (
  engineer_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  company_id STRING(36) NOT NULL,
  first_name STRING(50),
  last_name STRING(50),
  other_details STRING(255),
  CONSTRAINT fk_engineer_company FOREIGN KEY (company_id) REFERENCES Third_Party_Companies(company_id)
) PRIMARY KEY (engineer_id);

CREATE TABLE Engineer_Skills (
  engineer_id STRING(36) NOT NULL,
  skill_id STRING(36) NOT NULL,
  CONSTRAINT fk_engineer_skills_engineer FOREIGN KEY (engineer_id) REFERENCES Maintenance_Engineers(engineer_id),
  CONSTRAINT fk_engineer_skills_skill FOREIGN KEY (skill_id) REFERENCES Skills(skill_id)
) PRIMARY KEY (engineer_id, skill_id);

CREATE TABLE Fault_Log (
  fault_log_entry_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  asset_id STRING(36) NOT NULL,
  recorded_by_staff_id STRING(36) NOT NULL,
  fault_log_entry_datetime TIMESTAMP,
  fault_description STRING(255),
  other_fault_details STRING(255),
  CONSTRAINT fk_fault_log_asset FOREIGN KEY (asset_id) REFERENCES Assets(asset_id),
  CONSTRAINT fk_fault_log_recorded_by_staff FOREIGN KEY (recorded_by_staff_id) REFERENCES Staff(staff_id)
) PRIMARY KEY (fault_log_entry_id);

CREATE TABLE Engineer_Visits (
  engineer_visit_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  contact_staff_id STRING(36),
  engineer_id STRING(36) NOT NULL,
  fault_log_entry_id STRING(36) NOT NULL,
  fault_status STRING(10) NOT NULL,
  visit_start_datetime TIMESTAMP,
  visit_end_datetime TIMESTAMP,
  other_visit_details STRING(255),
  CONSTRAINT fk_engineer_visits_fault_log FOREIGN KEY (fault_log_entry_id) REFERENCES Fault_Log(fault_log_entry_id),
  CONSTRAINT fk_engineer_visits_engineer FOREIGN KEY (engineer_id) REFERENCES Maintenance_Engineers(engineer_id),
  CONSTRAINT fk_engineer_visits_contact_staff FOREIGN KEY (contact_staff_id) REFERENCES Staff(staff_id)
) PRIMARY KEY (engineer_visit_id);

CREATE TABLE Part_Faults (
  part_fault_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  part_id STRING(36) NOT NULL,
  fault_short_name STRING(20),
  fault_description STRING(255),
  other_fault_details STRING(255),
  CONSTRAINT fk_part_faults_part FOREIGN KEY (part_id) REFERENCES Parts(part_id)
) PRIMARY KEY (part_fault_id);

CREATE TABLE Fault_Log_Parts (
  fault_log_entry_id STRING(36) NOT NULL,
  part_fault_id STRING(36) NOT NULL,
  fault_status STRING(10) NOT NULL,
  CONSTRAINT fk_fault_log_parts_fault_log FOREIGN KEY (fault_log_entry_id) REFERENCES Fault_Log(fault_log_entry_id),
  CONSTRAINT fk_fault_log_parts_part_fault FOREIGN KEY (part_fault_id) REFERENCES Part_Faults(part_fault_id)
) PRIMARY KEY (fault_log_entry_id, part_fault_id);

CREATE TABLE Skills_Required_To_Fix (
  part_fault_id STRING(36) NOT NULL,
  skill_id STRING(36) NOT NULL,
  CONSTRAINT fk_skills_required_part_fault FOREIGN KEY (part_fault_id) REFERENCES Part_Faults(part_fault_id),
  CONSTRAINT fk_skills_required_skill FOREIGN KEY (skill_id) REFERENCES Skills(skill_id)
) PRIMARY KEY (part_fault_id, skill_id);