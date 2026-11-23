USE care_user;
INSERT INTO sys_user (mobile,password_hash,real_name,role) VALUES
('13800000001','123456','管理员A','admin'),
('13800000002','123456','家属B','family'),
('13800000003','123456','护士C','nurse');
INSERT INTO family_binding (family_user_id,elderly_id,relation) VALUES (2,1,'子女');

USE care_elderly;
INSERT INTO elderly_archive (real_name,gender,birth_date,mobile,address,emergency_contact,emergency_mobile,chronic_disease,care_level)
VALUES ('张爷爷',1,'1945-01-01','13900000001','幸福社区1号楼','张女士','13900000002','["高血压","脑卒中"]',4),
       ('李奶奶',2,'1948-03-02','13900000003','幸福社区2号楼','李先生','13900000004','["糖尿病"]',3);
INSERT INTO elderly_medication (elderly_id,nurse_id,drug_name,dosage,frequency,remind_times,start_date)
VALUES (1,1,'氨氯地平','5mg','每日1次','["08:00"]','2026-01-01');

USE care_nurse;
INSERT INTO nurse_profile (real_name,mobile,license_no,license_expire,skill_tags,service_area,current_load,score_avg)
VALUES ('王护士','13700000001','LIC001','2027-12-31','["压疮护理","更换胃管"]','幸福社区',0,4.8),
       ('赵护士','13700000002','LIC002','2026-12-31','["慢病护理","康复训练"]','幸福社区',0,4.9);
INSERT INTO training_record (nurse_id,course_name,exam_month,score,attempt_count,passed)
VALUES (1,'居家应急处置','2026-10',85,1,1);

USE care_vital;
INSERT INTO vital_record (elderly_id,heart_rate,systolic,diastolic,spo2,posture,source)
VALUES (1,78,130,82,97.0,1,'keyboard');

USE care_order;
INSERT INTO service_order (order_no,elderly_id,service_type,status,urgency,amount)
VALUES ('ORD20261024001',1,'压疮护理',0,1,199.00);

USE care_alert;
INSERT INTO adverse_event (event_type,description,anonymous,status)
VALUES ('fall','模拟匿名上报：老人跌倒未受伤',1,0);