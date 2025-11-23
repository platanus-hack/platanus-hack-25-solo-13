-- Drop learning plan tables in reverse order
DROP INDEX IF EXISTS idx_learning_plan_components_tipo;
DROP INDEX IF EXISTS idx_learning_plan_components_estado;
DROP INDEX IF EXISTS idx_learning_plan_components_plan_id;
DROP TABLE IF EXISTS learning_plan_components;

DROP INDEX IF EXISTS idx_learning_plans_estado;
DROP INDEX IF EXISTS idx_learning_plans_oa_bloom_objective_id;
DROP INDEX IF EXISTS idx_learning_plans_user_id;
DROP TABLE IF EXISTS learning_plans;
