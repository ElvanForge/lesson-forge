import { createClient } from "@supabase/supabase-js";
const supabaseUrl = "https://dpzkoqxfihjzontxhynp.supabase.co";
const supabaseAnonKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImRwemtvcXhmaWhqem9udHhoeW5wIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzE1ODc3MTQsImV4cCI6MjA4NzE2MzcxNH0.m57B9GoK-TV92MXbmmn9_DqnieYLF46v7aEAXXkkBnU";
const supabase = createClient(
  supabaseUrl,
  supabaseAnonKey
);
export {
  supabase as s
};
