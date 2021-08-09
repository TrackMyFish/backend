(this.webpackJsonpfrontend=this.webpackJsonpfrontend||[]).push([[0],{170:function(e,t){},221:function(e,t,a){},222:function(e,t,a){},223:function(e,t,a){},381:function(e,t,a){"use strict";a.r(t);var s=a(0),c=a.n(s),n=a(51),r=a.n(n),i=(a(199),a(221),a(222),a(170)),l=(a.n(i).a,a(34)),o=a(69),d=a(19),h=(a(223),a(3)),b=Object(l.a)((function(){return Object(h.jsx)(o.a,{children:Object(h.jsxs)(d.d,{children:[Object(h.jsx)(d.a,{exact:!0,from:"/",to:"/fish"}),Object(h.jsx)(d.b,{path:"/fish",component:g}),Object(h.jsx)(d.b,{path:"/tank",component:G})]})})})),u=a(87),j=a(9),m=a(41),p=a.n(m),O=a(42),f=new function e(){var t=this;Object(u.a)(this,e),this.fishState={fish:[],heartbeat:{fishbase:{status:""}},error:null},this.heartbeat=function(){p()({method:"get",headers:{"Content-Type":"application/json"},url:O.server.baseURL+"/heartbeat"}).then((function(e){Object(j.o)((function(){var a,s;t.fishState.heartbeat.fishbase.status=null===e||void 0===e||null===(a=e.data)||void 0===a||null===(s=a.fishbase)||void 0===s?void 0:s.status}))})).catch((function(e){var t;console.log(null===e||void 0===e||null===(t=e.response)||void 0===t?void 0:t.data)}))},this.listFish=function(){p()({method:"get",headers:{"Content-Type":"application/json"},url:O.server.baseURL+"/fish"}).then((function(e){Object(j.o)((function(){t.fishState.fish=e.data.fish}))})).catch((function(e){var t;console.log(null===e||void 0===e||null===(t=e.response)||void 0===t?void 0:t.data)}))},this.addFish=function(e){var a=e.genus,s=e.species,c=e.commonName,n=e.name,r=e.color,i=e.gender,l=e.purchaseDate,o=e.count;""===i&&(i="UNSPECIFIED"),p()({method:"post",headers:{"Content-Type":"application/json"},url:O.server.baseURL+"/fish",data:{genus:a,species:s,commonName:c,name:n,color:r,gender:i.toUpperCase(),purchaseDate:l,count:o}}).then((function(e){Object(j.o)((function(){var a;t.fishState.fish.push(null===e||void 0===e||null===(a=e.data)||void 0===a?void 0:a.fish)}))})).catch((function(e){Object(j.o)((function(){var a,s;console.log(e),t.fishState.error=null===e||void 0===e||null===(a=e.response)||void 0===a||null===(s=a.data)||void 0===s?void 0:s.message}))}))},this.removeFish=function(e){p()({method:"delete",headers:{"Content-Type":"application/json"},url:O.server.baseURL+"/fish/"+e}).then((function(){Object(j.o)((function(){t.fishState.fish=t.fishState.fish.filter((function(t){return t.id!==e}))}))})).catch((function(e){Object(j.o)((function(){var a,s;console.log(e),t.fishState.error=null===e||void 0===e||null===(a=e.response)||void 0===a||null===(s=a.data)||void 0===s?void 0:s.message}))}))},Object(j.l)(this)},v=c.a.createContext(f),x=new function e(){var t=this;Object(u.a)(this,e),this.tankState={stats:[],error:null},this.resetError=function(){t.tankState.error=null},this.listTankStatistics=function(){p()({method:"get",headers:{"Content-Type":"application/json"},url:O.server.baseURL+"/tank/statistics"}).then((function(e){Object(j.o)((function(){t.tankState.stats=e.data.tankStatistics}))})).catch((function(e){var t;console.log(null===e||void 0===e||null===(t=e.response)||void 0===t?void 0:t.data)}))},this.addTankStatistic=function(e){var a=e.testDate,s=e.ammonia,c=e.ph,n=e.nitrate,r=e.nitrite,i=e.gh,l=e.kh,o=e.phosphate;t.resetError(),p()({method:"post",headers:{"Content-Type":"application/json"},url:O.server.baseURL+"/tank/statistics",data:{testDate:a,ammonia:s,ph:c,nitrate:n,nitrite:r,gh:i,kh:l,phosphate:o}}).then((function(e){Object(j.o)((function(){var a;t.tankState.stats.push(null===e||void 0===e||null===(a=e.data)||void 0===a?void 0:a.tankStatistic)}))})).catch((function(e){Object(j.o)((function(){var a,s;console.log(e),t.tankState.error=null===e||void 0===e||null===(a=e.response)||void 0===a||null===(s=a.data)||void 0===s?void 0:s.message}))}))},this.removeTankStatistic=function(e){p()({method:"delete",headers:{"Content-Type":"application/json"},url:O.server.baseURL+"/tank/statistics/"+e}).then((function(){Object(j.o)((function(){t.tankState.stats=t.tankState.stats.filter((function(t){return t.id!==e}))}))})).catch((function(e){Object(j.o)((function(){var a,s;console.log(e),t.tankState.error=null===e||void 0===e||null===(a=e.response)||void 0===a||null===(s=a.data)||void 0===s?void 0:s.message}))}))},Object(j.l)(this)},N=c.a.createContext(x),g=Object(l.a)((function(){var e=Object(s.useContext)(v),t=e.fishState,a=e.heartbeat,c=e.listFish,n=e.addFish,r=e.removeFish;Object(s.useEffect)((function(){c(),a()}),[t]);var i=t.heartbeat.fishbase.status.toLowerCase();return Object(h.jsxs)("div",{children:[Object(h.jsx)(C,{}),Object(h.jsx)("div",{className:"mt-3 mb-3"}),i&&Object(h.jsxs)("div",{className:"container",children:[Object(h.jsxs)("div",{className:"alert alert-danger alert-dismissible",role:"alert",children:[Object(h.jsx)("h4",{className:"alert-heading",children:"Fishbase Down"}),Object(h.jsx)("p",{children:"Fishbase appears to be down, meaning certain functions such as the auto-population of the Ecosystem information will not be completed."}),Object(h.jsx)("hr",{}),Object(h.jsxs)("p",{className:"mb-0",children:["This is unlikely to be an issue with TrackMyFish and you should check the"," ",Object(h.jsx)("a",{href:"https://fishbase.ropensci.org/",className:"alert-link",children:"Fishbase API"})," ","for further details."]}),Object(h.jsx)("button",{type:"button",className:"btn-close","data-bs-dismiss":"alert","aria-label":"Close"})]}),Object(h.jsx)("div",{className:"mt-3 mb-3"})]}),Object(h.jsx)(y,{fish:t.fish,addFish:n,removeFish:r,fishError:t.error})]})})),C=function(){return Object(h.jsx)("nav",{className:"navbar navbar-expand-lg navbar-light",style:{backgroundColor:"#d7ecd9"},children:Object(h.jsxs)("div",{className:"container-fluid",children:[Object(h.jsx)("a",{className:"navbar-brand",href:"/home",children:"TrackMyFish"}),Object(h.jsx)("button",{className:"navbar-toggler",type:"button","data-bs-toggle":"collapse","data-bs-target":"#navbarSupportedContent","aria-controls":"navbarSupportedContent","aria-expanded":"false","aria-label":"Toggle navigation",children:Object(h.jsx)("span",{className:"navbar-toggler-icon"})}),Object(h.jsx)("div",{className:"collapse navbar-collapse",id:"navbarSupportedContent",children:Object(h.jsxs)("ul",{className:"navbar-nav me-auto mb-2 mb-lg-0",children:[Object(h.jsx)("li",{className:"nav-item",children:Object(h.jsx)(o.b,{className:"nav-link",to:"/fish",children:"Fish"})}),Object(h.jsx)("li",{className:"nav-item",children:Object(h.jsx)(o.b,{className:"nav-link",to:"/tank",children:"Tank"})}),Object(h.jsxs)("li",{className:"nav-item dropdown",children:[Object(h.jsx)("a",{className:"nav-link dropdown-toggle",href:"#",id:"navbarDropdown",role:"button","data-bs-toggle":"dropdown","aria-expanded":"false",children:"Useful Links"}),Object(h.jsx)("ul",{className:"dropdown-menu","aria-labelledby":"navbarDropdown",children:Object(h.jsx)("li",{children:Object(h.jsx)("a",{className:"dropdown-item",href:"https://www.fishbase.de/",children:"Fishbase"})})})]})]})})]})})},S=a(16),k=a(93),y=(a(138),Object(l.a)((function(e){var t=Object(s.useState)(""),a=Object(S.a)(t,2),c=a[0],n=a[1],r=Object(s.useState)(""),i=Object(S.a)(r,2),l=i[0],o=i[1],d=Object(s.useState)(""),b=Object(S.a)(d,2),u=b[0],j=b[1],m=Object(s.useState)(""),p=Object(S.a)(m,2),O=p[0],f=p[1],v=Object(s.useState)(""),x=Object(S.a)(v,2),N=x[0],g=x[1],C=Object(s.useState)(""),y=Object(S.a)(C,2),F=y[0],H=y[1],E=Object(s.useState)(""),T=Object(S.a)(E,2),w=T[0],D=T[1],P=Object(s.useState)(0),R=Object(S.a)(P,2),G=R[0],L=R[1],U=Object(s.useState)(""),I=Object(S.a)(U,2),A=I[0],K=I[1],J=e.fishError||A;return Object(h.jsxs)("div",{className:"container-fluid",children:[Object(h.jsx)("h3",{className:"text-center",children:"Fish"}),Object(h.jsx)(k.a,{data:e.fish,sortable:!0,pageSize:e.fish.length,columns:[{Header:"ID",accessor:"id",headerClassName:"text-start"},{Header:"Genus",accessor:"genus",headerClassName:"text-start"},{Header:"Species",accessor:"species",headerClassName:"text-start"},{Header:"Common Name",accessor:"commonName",headerClassName:"text-start"},{Header:"Name",accessor:"name",headerClassName:"text-start"},{Header:"Color",accessor:"color",headerClassName:"text-start"},{Header:"Gender",accessor:"gender",headerClassName:"text-start",Cell:function(e){var t=e.value;return"UNSPECIFIED"===t?"":t}},{Header:"Purchase Date",accessor:"purchaseDate",headerClassName:"text-start"},{Header:"Ecosystem",columns:[{Header:"Name",accessor:"ecosystemName",headerClassName:"text-start"},{Header:"Type",accessor:"ecosystemType",headerClassName:"text-start"},{Header:"Location",accessor:"ecosystemLocation",headerClassName:"text-start"},{Header:"Salinity",accessor:"salinity",headerClassName:"text-start"},{Header:"Climate",accessor:"climate",headerClassName:"text-start"}]},{Header:"Count",accessor:"count",headerClassName:"text-start"},{className:"text-center",Cell:function(t){return Object(h.jsx)("button",{className:"btn btn-danger ml-2",onClick:function(){return e.removeFish(t.original.id)},children:"Remove"})}}]}),Object(h.jsx)("div",{className:"accordion",id:"create-fish-accordion",children:Object(h.jsxs)("div",{className:"accordion-item",children:[Object(h.jsx)("h2",{className:"accordion-header",id:"fish-heading",children:Object(h.jsx)("button",{className:"accordion-button collapsed",type:"button","data-bs-toggle":"collapse","data-bs-target":"#collapse-fish","aria-expanded":"false","aria-controls":"collapse-fish",children:"Add Fish"})}),Object(h.jsx)("div",{id:"collapse-fish",className:"accordion-collapse collapse","aria-labelledby":"fish-heading","data-bs-parent":"#create-fish-accordion",children:Object(h.jsxs)("div",{className:"accordion-body",children:[Object(h.jsxs)("form",{children:[Object(h.jsxs)("div",{className:"container-fluid ps-0 pe-0",children:[Object(h.jsxs)("div",{className:"row",children:[Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputGenus",className:"form-label",children:"Genus"}),Object(h.jsx)("input",{type:"text",className:"form-control",id:"inputGenus",onChange:function(e){return K("")||n(e.target.value)},value:c})]})}),Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputSpecies",className:"form-label",children:"Species"}),Object(h.jsx)("input",{type:"text",className:"form-control",id:"inputSpecies",onChange:function(e){return K("")||o(e.target.value)},value:l})]})})]}),Object(h.jsxs)("div",{className:"row",children:[Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputCommonName",className:"form-label",children:"Common Name"}),Object(h.jsx)("input",{type:"text",className:"form-control",id:"inputCommonName",onChange:function(e){return K("")||j(e.target.value)},value:u})]})}),Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputName",className:"form-label",children:"Name"}),Object(h.jsx)("input",{type:"text",className:"form-control",id:"inputName",onChange:function(e){return K("")||f(e.target.value)},value:O})]})})]}),Object(h.jsxs)("div",{className:"row",children:[Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputColor",className:"form-label",children:"Color"}),Object(h.jsx)("input",{type:"text",className:"form-control",id:"inputColor",onChange:function(e){return K("")||g(e.target.value)},value:N})]})}),Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputGender",className:"form-label",children:"Gender"}),Object(h.jsx)("input",{type:"text",className:"form-control",id:"inputGender",onChange:function(e){return K("")||H(e.target.value)},value:F})]})}),Object(h.jsxs)("div",{className:"row",children:[Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputPurchaseDate",className:"form-label",children:"Purchase Date"}),Object(h.jsx)("input",{type:"date",className:"form-control",id:"inputPurchaseDate",onChange:function(e){return K("")||D(e.target.value)},value:w})]})}),Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputCount",className:"form-label",children:"Count"}),Object(h.jsx)("input",{type:"number",className:"form-control",id:"inputCount",min:"0",onChange:function(e){return K("")||L(e.target.value)},value:G})]})})]})]})]}),Object(h.jsx)("button",{type:"submit",className:"btn btn-primary",onClick:function(t){t.preventDefault();var a=F.toUpperCase();e.addFish({genus:c,species:l,commonName:u,name:O,color:N,gender:F,uppercaseGender:a,purchaseDate:w,count:G})},children:"Submit"})]}),J&&Object(h.jsx)("div",{className:"alert alert-danger mt-3",role:"alert",children:J})]})})]})})]})}))),F=a(383),H=a(384),E=a(388),T=a(187),w=a(188),D=a(78),P=a(75),R=a(190),G=Object(l.a)((function(){var e=Object(s.useContext)(N),t=e.tankState,a=e.resetError,c=e.listTankStatistics,n=e.addTankStatistic,r=e.removeTankStatistic;return Object(s.useEffect)((function(){c()}),[t]),Object(h.jsxs)("div",{children:[Object(h.jsx)(C,{}),Object(h.jsx)("div",{className:"mt-3 mb-3"}),Object(h.jsx)(L,{tankStatistics:t.stats,addTankStatistic:n,removeTankStatistic:r,tankStatisticsError:t.error,resetError:a}),Object(h.jsx)("div",{className:"mt-3 mb-3"}),Object(h.jsx)(F.a,{width:"100%",height:"100%",children:Object(h.jsxs)(H.a,{data:t.stats,margin:{top:5,right:20,left:10,bottom:5},children:[Object(h.jsx)(E.a,{strokeDasharray:"3 3"}),Object(h.jsx)(T.a,{dataKey:"testDate",padding:{left:30,right:30}}),Object(h.jsx)(w.a,{}),Object(h.jsx)(D.a,{}),Object(h.jsx)(P.a,{}),Object(h.jsx)(R.a,{connectNulls:!0,type:"monotone",dataKey:"ammonia",stroke:"#ff7300"}),Object(h.jsx)(R.a,{connectNulls:!0,type:"monotone",dataKey:"ph",stroke:"#387908"})]})})]})})),L=Object(l.a)((function(e){var t=Object(s.useState)(""),a=Object(S.a)(t,2),c=a[0],n=a[1],r=Object(s.useRef)(null),i=Object(s.useRef)(null),l=Object(s.useRef)(null),o=Object(s.useRef)(null),d=Object(s.useRef)(null),b=Object(s.useRef)(null),u=Object(s.useRef)(null);return Object(h.jsxs)("div",{className:"container-fluid",children:[Object(h.jsx)("h3",{className:"text-center",children:"Tank Statistics"}),Object(h.jsx)(k.a,{data:e.tankStatistics,sortable:!0,pageSize:e.tankStatistics.length,columns:[{Header:"ID",accessor:"id",headerClassName:"text-start"},{Header:"Test Date",accessor:"testDate",headerClassName:"text-start"},{Header:"Ammonia",accessor:"ammonia",headerClassName:"text-start"},{Header:"pH",accessor:"ph",headerClassName:"text-start"},{Header:"Nitite",accessor:"nitrite",headerClassName:"text-start"},{Header:"Nitrate",accessor:"nitrate",headerClassName:"text-start"},{Header:"GH",accessor:"gh",headerClassName:"text-start"},{Header:"KH",accessor:"kh",headerClassName:"text-start"},{Header:"Phosphate",accessor:"phosphate",headerClassName:"text-start"},{className:"text-center",Cell:function(t){return Object(h.jsx)("button",{className:"btn btn-danger ml-2",onClick:function(){return e.removeTankStatistic(t.original.id)},children:"Remove"})}}]}),Object(h.jsx)("div",{className:"accordion",id:"create-tank-statistic-accordion",children:Object(h.jsxs)("div",{className:"accordion-item",children:[Object(h.jsx)("h2",{className:"accordion-header",id:"fish-heading",children:Object(h.jsx)("button",{className:"accordion-button collapsed",type:"button","data-bs-toggle":"collapse","data-bs-target":"#collapse-tank-statistic","aria-expanded":"false","aria-controls":"collapse-tank-statistic",children:"Add Tank Statistics"})}),Object(h.jsx)("div",{id:"collapse-tank-statistic",className:"accordion-collapse collapse","aria-labelledby":"fish-heading","data-bs-parent":"#create-tank-statistic-accordion",children:Object(h.jsxs)("div",{className:"accordion-body",children:[Object(h.jsxs)("form",{children:[Object(h.jsxs)("div",{className:"container-fluid ps-0 pe-0",children:[Object(h.jsxs)("div",{className:"row",children:[Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputTestDate",className:"form-label",children:"Test Date"}),Object(h.jsx)("input",{type:"date",className:"form-control",id:"inputTestDate",onChange:function(e){return n(e.target.value)},value:c})]})}),Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputPH",className:"form-label",children:"pH"}),Object(h.jsx)("input",{type:"number",step:"0.01",className:"form-control",id:"inputPH",placeholder:"pH",ref:r,onChange:function(){return e.tankStatisticsError&&e.resetError()}})]})})]}),Object(h.jsxs)("div",{className:"row",children:[Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputAmmonia",className:"form-label",children:"Ammonia"}),Object(h.jsx)("input",{type:"number",step:"0.01",className:"form-control",id:"inputAmmonia",placeholder:"Ammonia (ppm)",ref:i,onChange:function(){return e.tankStatisticsError&&e.resetError()}})]})}),Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputNitrite",className:"form-label",children:"Nitrite (NO2)"}),Object(h.jsx)("input",{type:"number",step:"0.01",className:"form-control",id:"inputNitrite",placeholder:"Nitrite (ppm)",ref:o,onChange:function(){return e.tankStatisticsError&&e.resetError()}})]})})]}),Object(h.jsxs)("div",{className:"row",children:[Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputNitrate",className:"form-label",children:"Nitrate (NO3)"}),Object(h.jsx)("input",{type:"number",step:"0.01",className:"form-control",id:"inputNitrate",placeholder:"Nitrate (ppm)",ref:l,onChange:function(){return e.tankStatisticsError&&e.resetError()}})]})}),Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputPhosphate",className:"form-label",children:"Phosphate (PO4)"}),Object(h.jsx)("input",{type:"number",step:"0.01",className:"form-control",id:"inputPhosphate",placeholder:"Phosphate (ppm)",ref:u,onChange:function(){return e.tankStatisticsError&&e.resetError()}})]})})]}),Object(h.jsxs)("div",{className:"row",children:[Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputGH",className:"form-label",children:"GH"}),Object(h.jsx)("input",{type:"number",step:"0.01",className:"form-control",id:"inputGH",placeholder:"GH",ref:d,onChange:function(){return e.tankStatisticsError&&e.resetError()}})]})}),Object(h.jsx)("div",{className:"col",children:Object(h.jsxs)("div",{className:"mb-3",children:[Object(h.jsx)("label",{htmlFor:"inputKH",className:"form-label",children:"KH"}),Object(h.jsx)("input",{type:"number",step:"0.01",className:"form-control",id:"inputKH",placeholder:"KH",ref:b,onChange:function(){return e.tankStatisticsError&&e.resetError()}})]})})]})]}),Object(h.jsx)("button",{type:"submit",className:"btn btn-primary",onClick:function(t){t.preventDefault(),e.addTankStatistic({testDate:c,ammonia:i.current.value?i.current.value:null,ph:r.current.value?r.current.value:null,nitrate:l.current.value?l.current.value:null,nitrite:o.current.value?o.current.value:null,gh:d.current.value?d.current.value:null,kh:b.current.value?b.current.value:null,phosphate:u.current.value?u.current.value:null})},children:"Submit"})]}),e.tankStatisticsError&&Object(h.jsx)("div",{className:"alert alert-danger mt-3",role:"alert",children:e.tankStatisticsError})]})})]})})]})})),U=function(e){e&&e instanceof Function&&a.e(3).then(a.bind(null,390)).then((function(t){var a=t.getCLS,s=t.getFID,c=t.getFCP,n=t.getLCP,r=t.getTTFB;a(e),s(e),c(e),n(e),r(e)}))};r.a.render(Object(h.jsx)(b,{}),document.getElementById("root")),U()},42:function(e){e.exports=JSON.parse('{"server":{"host":"localhost","port":8443,"baseURL":"/api/v1alpha1"}}')}},[[381,1,2]]]);
//# sourceMappingURL=main.a6066072.chunk.js.map