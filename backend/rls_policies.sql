-- Enable RLS on all tables
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE generations ENABLE ROW LEVEL SECURITY;

-- 1. Policies for 'users' table
-- Users can only read their own profile
CREATE POLICY "Users can view own profile" 
ON users FOR SELECT 
TO authenticated 
USING (auth.uid() = id);

-- CRITICAL: Users cannot INSERT or UPDATE their own balance.
-- This is handled by the Go backend (service_role) or Triggers.

-- 2. Trigger for New User Capture
-- Automatically creates a profile and grants 10 starting credits on signup
CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS trigger AS $$
BEGIN
  INSERT INTO public.users (id, email, credit_balance)
  VALUES (new.id, new.email, 10);
  RETURN new;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Trigger the function every time a user is created in auth.users
CREATE OR REPLACE TRIGGER on_auth_user_created
  AFTER INSERT ON auth.users
  FOR EACH ROW EXECUTE FUNCTION public.handle_new_user();

-- 3. Policies for 'generations' table
-- Users can only see their own generations
CREATE POLICY "Users can view own generations" 
ON generations FOR SELECT 
TO authenticated 
USING (auth.uid() = user_id);

-- Users can insert their own generations
CREATE POLICY "Users can insert own generations" 
ON generations FOR INSERT 
TO authenticated 
WITH CHECK (auth.uid() = user_id);

-- 3. Policies for 'transactions' table
-- Users can only see their own transactions
CREATE POLICY "Users can view own transactions" 
ON transactions FOR SELECT 
TO authenticated 
USING (auth.uid() = user_id);
