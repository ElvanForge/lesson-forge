import { e as escape_html, a2 as attr, $ as attr_class, a1 as stringify } from "../../../chunks/index.js";
import { s as supabase } from "../../../chunks/supabase.js";
import { B as Button } from "../../../chunks/Button.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let email = "";
    let password = "";
    let isRegistering = false;
    let isLoading = false;
    let message = { text: "", type: "" };
    async function handleAuth() {
      isLoading = true;
      message = { text: "", type: "" };
      try {
        if (isRegistering) ;
        else {
          const { error } = await supabase.auth.signInWithPassword({ email, password });
          if (error) throw error;
          window.location.href = "/";
        }
      } catch (error) {
        message = { text: error.message, type: "error" };
      } finally {
        isLoading = false;
      }
    }
    $$renderer2.push(`<div class="min-h-screen bg-slate-50 flex items-center justify-center p-4"><div class="max-w-md w-full animate-fade-in"><div class="text-center mb-8"><h1 class="text-3xl font-bold text-primary mb-2">Lesson Forge</h1> <p class="text-slate-500 font-medium">Helping Forge Future Minds</p></div> <div class="card-premium p-8"><h2 class="text-2xl font-bold text-slate-800 mb-6 text-center">${escape_html("Welcome Back")}</h2> `);
    if (message.text) {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<div${attr_class(`mb-6 p-4 rounded-xl text-sm font-medium ${stringify(message.type === "error" ? "bg-red-50 text-red-600 border border-red-100" : "bg-emerald-50 text-emerald-600 border border-emerald-100")} animate-fade-in`)}>${escape_html(message.text)}</div>`);
    } else {
      $$renderer2.push("<!--[!-->");
    }
    $$renderer2.push(`<!--]--> <form class="space-y-4"><div><label for="email" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Email Address</label> <input type="email" id="email"${attr("value", email)} required="" class="w-full px-4 py-3 bg-slate-50 border-2 border-slate-100 rounded-xl focus:ring-4 focus:ring-primary/10 focus:border-primary transition-all text-slate-800 placeholder-slate-400" placeholder="teacher@example.com"/></div> <div><label for="password" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Password</label> <input type="password" id="password"${attr("value", password)} required="" class="w-full px-4 py-3 bg-slate-50 border-2 border-slate-100 rounded-xl focus:ring-4 focus:ring-primary/10 focus:border-primary transition-all text-slate-800 placeholder-slate-400" placeholder="••••••••"/></div> <div class="pt-2">`);
    Button($$renderer2, {
      text: "Sign In",
      isLoading,
      onclick: handleAuth
    });
    $$renderer2.push(`<!----></div></form> <div class="mt-8 relative"><div class="absolute inset-0 flex items-center"><div class="w-full border-t border-slate-200"></div></div> <div class="relative flex justify-center text-xs uppercase"><span class="bg-white px-4 text-slate-400 font-bold tracking-widest">Or continue with</span></div></div> <div class="mt-6"><button class="w-full flex items-center justify-center gap-3 px-4 py-3 bg-white border-2 border-slate-100 rounded-xl hover:border-accent hover:bg-slate-50 transition-all group"><svg class="w-5 h-5 group-hover:scale-110 transition-transform" viewBox="0 0 24 24"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"></path><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"></path><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"></path><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"></path></svg> <span class="text-sm font-bold text-slate-700">Google</span></button></div> <div class="mt-8 text-center"><button class="text-sm font-bold text-primary hover:text-primary/80 transition-colors">${escape_html("Don't have an account? Register Now")}</button></div></div> <div class="mt-8 text-center"><a href="/" class="text-sm font-medium text-slate-400 hover:text-slate-600 transition-colors flex items-center justify-center gap-2"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg> Back to Home</a></div></div></div>`);
  });
}
export {
  _page as default
};
