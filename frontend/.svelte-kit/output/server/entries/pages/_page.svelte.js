import { e as escape_html, $ as attr_class, a0 as attr, a1 as ensure_array_like, a2 as stringify } from "../../chunks/index.js";
import "clsx";
import { B as Button } from "../../chunks/Button.js";
import "../../chunks/supabase.js";
function Header($$renderer, $$props) {
  let { title, credits } = $$props;
  $$renderer.push(`<header class="glass sticky top-0 z-50 border-b border-white/40"><div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-20 flex items-center justify-between"><a href="/" class="flex items-center gap-3 hover:opacity-80 transition-opacity"><div class="w-10 h-10 bg-primary rounded-xl flex items-center justify-center shadow-lg shadow-primary/20"><svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg></div> <h1 class="text-2xl font-extrabold tracking-tight text-primary">${escape_html(title)}</h1></a> <div class="flex items-center gap-6">`);
  {
    $$renderer.push("<!--[!-->");
  }
  $$renderer.push(`<!--]--> <div class="flex items-center gap-3"><span class="flex items-center gap-1.5 text-sm font-bold bg-primary/10 text-primary px-4 py-2 rounded-full border border-primary/10"><span class="w-2 h-2 rounded-full bg-primary animate-pulse"></span> ${escape_html(credits)} Credits</span> `);
  {
    $$renderer.push("<!--[!-->");
    $$renderer.push(`<a href="/login" class="text-sm font-bold text-primary hover:text-primary/80 transition-all">Sign In</a>`);
  }
  $$renderer.push(`<!--]--></div></div></div></header>`);
}
function EmptyState($$renderer, $$props) {
  let { message } = $$props;
  $$renderer.push(`<div class="text-center py-20 animate-fade-in"><div class="inline-flex items-center justify-center w-20 h-20 rounded-3xl bg-slate-50 mb-6 text-slate-300"><svg class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg></div> <h3 class="text-lg font-bold text-slate-800">No Generations Yet</h3> <p class="mt-2 text-slate-500 max-w-xs mx-auto">${escape_html(message)}</p></div>`);
}
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let credits = 0;
    let prompt = "";
    let history = [];
    async function handleGenerate() {
      return;
    }
    $$renderer2.push(`<div class="min-h-screen bg-[#F8FAFC]">`);
    Header($$renderer2, {
      credits,
      title: "Lesson Forge"
    });
    $$renderer2.push(`<!----> <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12"><div class="grid grid-cols-1 lg:grid-cols-12 gap-12"><div class="lg:col-span-8 space-y-8"><div class="bg-white rounded-3xl p-8 shadow-sm border border-slate-100"><h2 class="text-2xl font-bold text-slate-900 mb-6">Lesson Forge</h2> <div class="space-y-6"><div class="flex p-1 bg-slate-100 rounded-2xl w-fit"><button${attr_class(`px-6 py-2 rounded-xl text-sm font-medium transition-all ${stringify(
      "bg-white text-primary shadow-sm"
    )}`)}>Lesson Plan</button> <button${attr_class(`px-6 py-2 rounded-xl text-sm font-medium transition-all ${stringify("text-slate-500 hover:text-slate-700")}`)}>Presentation</button></div> <div class="relative"><textarea${attr(
      "placeholder",
      "e.g., A 45-minute ESL lesson..."
    )} class="w-full h-40 p-6 bg-slate-50 border-none rounded-3xl focus:ring-2 focus:ring-primary/20 transition-all resize-none text-slate-700 placeholder:text-slate-400">`);
    const $$body = escape_html(prompt);
    if ($$body) {
      $$renderer2.push(`${$$body}`);
    }
    $$renderer2.push(`</textarea></div> <div class="flex items-center justify-between bg-slate-50 p-4 rounded-2xl"><div class="flex items-center gap-3"><div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary"><svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg></div> <div><p class="text-sm font-bold text-slate-900">Cost: 1 Credits</p> <p class="text-xs text-slate-500">Current Balance: ${escape_html(credits)}</p></div></div> `);
    Button($$renderer2, {
      onclick: handleGenerate,
      disabled: !prompt,
      variant: "primary",
      text: "Forge Content"
    });
    $$renderer2.push(`<!----></div> `);
    {
      $$renderer2.push("<!--[!-->");
    }
    $$renderer2.push(`<!--]--></div></div> <div class="space-y-4"><h3 class="text-xl font-bold text-slate-900">Your Forge History</h3> `);
    if (history.length === 0) {
      $$renderer2.push("<!--[-->");
      EmptyState($$renderer2, {
        message: "Your generated lessons and presentations will appear here."
      });
    } else {
      $$renderer2.push("<!--[!-->");
      $$renderer2.push(`<!--[-->`);
      const each_array = ensure_array_like(history);
      for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
        let item = each_array[$$index];
        $$renderer2.push(`<div class="bg-white p-6 rounded-3xl border border-slate-100 flex items-center justify-between hover:shadow-md transition-all group"><div class="flex items-center gap-4"><div class="w-12 h-12 rounded-2xl bg-slate-50 flex items-center justify-center text-slate-400 group-hover:bg-primary/5 group-hover:text-primary transition-colors"><svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg></div> <div><p class="font-bold text-slate-900 line-clamp-1">${escape_html(item.prompt)}</p> <p class="text-xs text-slate-400">${escape_html(new Date(item.created_at).toLocaleDateString())}</p></div></div> <a${attr("href", item.file_path)} aria-label="Download generated file" target="_blank" class="p-3 rounded-xl bg-slate-50 text-slate-400 hover:bg-primary hover:text-white transition-all"><svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg></a></div>`);
      }
      $$renderer2.push(`<!--]-->`);
    }
    $$renderer2.push(`<!--]--></div></div> <div class="lg:col-span-4"><div class="p-8 bg-primary rounded-3xl text-white shadow-xl sticky top-8"><h3 class="text-xl font-bold mb-2">Fuel Your Forge</h3> <div class="space-y-4 mt-8"><button class="w-full bg-white text-primary font-bold py-4 rounded-2xl shadow-md hover:-translate-y-1 transition-all">10 Credits | $9.99</button> <button class="w-full bg-accent text-primary font-extrabold py-5 rounded-2xl shadow-lg relative hover:-translate-y-1 transition-all"><span class="absolute -top-3 left-1/2 -translate-x-1/2 bg-white text-primary text-[10px] px-3 py-1 rounded-full uppercase tracking-wider font-black">Best Value</span> 50 Credits | $39.99</button></div></div></div></div></main></div>`);
  });
}
export {
  _page as default
};
