const o="";async function r(n){const t=await fetch(o+n);if(!t.ok)throw new Error(`HTTP ${t.status}`);return t.json()}export{r as g};
