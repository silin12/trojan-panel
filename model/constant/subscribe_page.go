package constant

// SubscribePageTemplate 浏览器访问订阅短链时展示的信息页
const SubscribePageTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{SYSTEM_NAME}} · 订阅信息</title>
<style>
:root{
  --bg:#101419;--panel:#161c23;--line:#232b34;--text:#e6edf3;--muted:#8b98a5;
  --accent:#4c9aff;--ok:#38b48b;--warn:#e8a33d;
}
*{box-sizing:border-box}
body{margin:0;min-height:100vh;background:var(--bg);color:var(--text);
font-family:-apple-system,"Segoe UI","PingFang SC","Microsoft YaHei",sans-serif;
display:flex;align-items:flex-start;justify-content:center;padding:40px 18px 60px}
.wrap{width:100%;max-width:760px}
.brand{display:flex;align-items:center;gap:10px;margin-bottom:22px}
.dot{width:10px;height:10px;border-radius:50%;background:var(--accent)}
.brand h1{font-size:15px;font-weight:600;letter-spacing:.04em;margin:0;color:var(--muted)}
.card{background:var(--panel);border:1px solid var(--line);border-radius:14px;padding:26px}
.head{display:flex;flex-wrap:wrap;align-items:baseline;gap:12px;margin-bottom:6px}
.head h2{margin:0;font-size:26px;font-weight:650;word-break:break-all}
.tag{font-size:12px;color:var(--muted);border:1px solid var(--line);border-radius:999px;padding:3px 10px}
.grid{display:grid;grid-template-columns:repeat(3,1fr);gap:14px;margin:22px 0 20px}
.metric{background:#1b222b;border:1px solid var(--line);border-radius:10px;padding:14px}
.metric .k{font-size:12px;color:var(--muted);margin-bottom:6px}
.metric .v{font-size:19px;font-weight:600;font-variant-numeric:tabular-nums;word-break:break-all}
.bar{height:8px;background:#1b222b;border:1px solid var(--line);border-radius:999px;overflow:hidden}
.bar > i{display:block;height:100%;background:linear-gradient(90deg,var(--ok),var(--accent));width:0}
.barNote{display:flex;justify-content:space-between;font-size:12px;color:var(--muted);margin-top:8px}
.section{margin-top:26px}
.section h3{font-size:13px;font-weight:600;color:var(--muted);margin:0 0 10px;letter-spacing:.03em}
.row{display:flex;gap:10px;align-items:center;background:#1b222b;border:1px solid var(--line);
border-radius:10px;padding:10px 12px;margin-bottom:10px}
.row .name{flex:0 0 92px;font-size:13px;color:var(--muted)}
.row .link{flex:1;min-width:0;font-family:ui-monospace,SFMono-Regular,Menlo,monospace;
font-size:12px;color:#b8c4d0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.row button{flex:0 0 auto;background:#22303d;color:var(--text);border:1px solid var(--line);
border-radius:8px;padding:7px 14px;font-size:13px;cursor:pointer}
.row button:hover{border-color:var(--accent);color:#fff}
.tip{font-size:12px;line-height:1.7;color:var(--muted);margin-top:18px}
.foot{margin-top:22px;font-size:12px;color:#5d6b7a;text-align:center}
@media (max-width:560px){
  .grid{grid-template-columns:1fr}
  .row{flex-wrap:wrap}
  .row .name{flex:0 0 auto}
  .row .link{flex:1 1 100%;order:3}
  .row button{flex:1 1 100%}
}
</style>
</head>
<body>
<div class="wrap">
  <div class="brand"><span class="dot"></span><h1>{{SYSTEM_NAME}} / SUBSCRIBE</h1></div>
  <div class="card">
    <div class="head">
      <h2>{{USERNAME}}</h2>
      <span class="tag">到期 {{EXPIRE}}</span>
    </div>
    <div class="grid">
      <div class="metric"><div class="k">剩余流量</div><div class="v">{{REMAIN}}</div></div>
      <div class="metric"><div class="k">总流量</div><div class="v">{{TOTAL}}</div></div>
      <div class="metric"><div class="k">已使用</div><div class="v">{{USED}}</div></div>
    </div>
    <div class="bar"><i style="width:{{PERCENT}}%"></i></div>
    <div class="barNote"><span>已使用 {{PERCENT}}%</span><span>{{USED}} / {{TOTAL}}</span></div>

    <div class="section">
      <h3>订阅地址（复制到客户端使用）</h3>
      <div class="row">
        <span class="name">Clash / Mihomo</span>
        <span class="link" id="u-clash">/api/s/{{CODE}}?format=clash</span>
        <button type="button" data-target="u-clash">复制</button>
      </div>
      <div class="row">
        <span class="name">v2rayN / 小火箭</span>
        <span class="link" id="u-v2ray">/api/s/{{CODE}}?format=v2ray</span>
        <button type="button" data-target="u-v2ray">复制</button>
      </div>
      <div class="row">
        <span class="name">sing-box</span>
        <span class="link" id="u-singbox">/api/s/{{CODE}}?format=singbox</span>
        <button type="button" data-target="u-singbox">复制</button>
      </div>
      <div class="row">
        <span class="name">通用 Base64</span>
        <span class="link" id="u-base64">/api/s/{{CODE}}?format=base64</span>
        <button type="button" data-target="u-base64">复制</button>
      </div>
      <div class="tip">
        直接粘贴 /api/s/{{CODE}} 这个短地址也可以：客户端会自动按自己的类型拿到对应格式，
        浏览器打开则显示本页面。
      </div>
    </div>
  </div>
  <div class="foot">{{SYSTEM_NAME}} · 订阅信息基于当前流量与到期时间生成</div>
</div>
<script>
(function(){
  function abs(id){
    var el=document.getElementById(id);
    if(!el){return "";}
    return new URL(el.textContent.trim(), location.origin).href;
  }
  function copy(text,cb){
    if(navigator.clipboard && window.isSecureContext){
      navigator.clipboard.writeText(text).then(cb, function(){ fallback(text,cb); });
    }else{
      fallback(text,cb);
    }
  }
  function fallback(text,cb){
    var ta=document.createElement("textarea");
    ta.value=text;ta.setAttribute("readonly","");
    ta.style.position="fixed";ta.style.top="-1000px";
    document.body.appendChild(ta);ta.select();
    try{document.execCommand("copy");}catch(e){}
    document.body.removeChild(ta);
    cb();
  }
  var buttons=document.querySelectorAll("button[data-target]");
  for(var i=0;i<buttons.length;i++){
    buttons[i].addEventListener("click",function(){
      var self=this;
      var link=abs(self.getAttribute("data-target"));
      copy(link,function(){
        var old=self.textContent;
        self.textContent="已复制";
        setTimeout(function(){self.textContent=old;},1500);
      });
    });
  }
  var links=document.querySelectorAll(".link");
  for(var j=0;j<links.length;j++){
    links[j].title=abs(links[j].id);
  }
})();
</script>
</body>
</html>`
