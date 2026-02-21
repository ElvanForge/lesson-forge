import { createClient } from "@supabase/supabase-js";
const supabaseUrl = "https://dpzkoqxfihjzontxhynp.supabase.co";
const supabaseAnonKey = "sb_publishable_1J79QPwHEmyQqjN0rIdTYw_OIXKEhF4";
const supabase = createClient(
  supabaseUrl,
  supabaseAnonKey
);
export {
  supabase as s
};
