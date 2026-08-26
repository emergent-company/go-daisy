package handler

const nxRevenueStatsConfig = `{
  "chart":{"height":288,"type":"bar","stacked":true,"background":"transparent","toolbar":{"show":false}},
  "plotOptions":{"bar":{"borderRadius":8,"borderRadiusApplication":"end","borderRadiusWhenStacked":"last","colors":{"backgroundBarColors":["rgba(150,150,150,0.07)"],"backgroundBarRadius":8},"columnWidth":"45%","barHeight":"100%"}},
  "dataLabels":{"enabled":false},
  "colors":["#ff8b4b","#6c74f8"],
  "legend":{"show":true,"horizontalAlign":"center","offsetX":0,"offsetY":6},
  "series":[{"name":"Orders","data":[10,12,14,16,18,20,14,16,24,12]},{"name":"Revenue","data":[15,24,21,28,30,40,22,32,48,20]}],
  "xaxis":{"categories":["2016","2017","2018","2019","2020","2021","2022","2023","2024","2025"],"axisBorder":{"show":false},"axisTicks":{"show":false}},
  "yaxis":{"axisBorder":{"show":false},"axisTicks":{"show":false},"labels":{"show":false}},
  "tooltip":{"enabled":true,"shared":true,"intersect":false},
  "grid":{"show":false},
  "responsive":[{"breakpoint":450,"options":{"plotOptions":{"bar":{"borderRadius":4}}}}]
}`

const nxCustomerAcqConfig = `{
  "chart":{"height":356,"type":"line","sparkline":{"enabled":false},"toolbar":{"show":false},"zoom":{"enabled":false},"background":"transparent"},
  "forecastDataPoints":{"count":2,"dashArray":[6,4]},
  "grid":{"show":false},
  "yaxis":{"show":false,"min":125,"max":181},
  "stroke":{"curve":"stepline","width":[2,1.5]},
  "colors":["#167bff","rgba(150,150,150,0.3)"],
  "series":[{"name":"Customer","data":[144,150,146,154,150,155,160,155,140,155,160,180,170,165,165]},{"name":"Advertise","data":[140,142,142,140,146,148,150,136,130,133,145,148,158,150,150]}],
  "xaxis":{"categories":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15]}
}`

const nxCRMSaleMetricsConfig = `{
  "chart":{"height":323,"type":"bar","background":"transparent","toolbar":{"show":false}},
  "plotOptions":{"bar":{"borderRadius":2,"columnWidth":"50%","barHeight":"100%"}},
  "stroke":{"show":true,"width":2,"colors":["transparent"]},
  "dataLabels":{"enabled":false},
  "colors":["#ff8b4b","#6c74f8"],
  "legend":{"show":true,"horizontalAlign":"center","offsetY":6},
  "series":[{"name":"Customer","data":[1175,1734,2239,2741,1823,2154,1013,2794,1834,3273]},{"name":"Acquisition","data":[1803,2175,2882,2486,3755,1888,3154,4345,2683,2891]}],
  "xaxis":{"categories":["2016","2017","2018","2019","2020","2021","2022","2023","2024","2025"],"axisBorder":{"show":false},"axisTicks":{"show":false}},
  "yaxis":{"axisBorder":{"show":false},"axisTicks":{"show":false}},
  "tooltip":{"enabled":true,"shared":true,"intersect":false},
  "grid":{"show":false},
  "responsive":[{"breakpoint":450,"options":{"plotOptions":{"bar":{"borderRadius":2}}}}]
}`

const nxCRMGoalStatusConfig = `{
  "series":[76],
  "chart":{"height":258,"type":"radialBar","offsetY":-20,"background":"transparent","sparkline":{"enabled":true}},
  "stroke":{"lineCap":"round"},
  "colors":["#167bff"],
  "plotOptions":{"radialBar":{"startAngle":-90,"endAngle":90,"track":{"background":"var(--color-base-200)","strokeWidth":"75%","margin":8,"dropShadow":{"enabled":false}}}}
}`

const nxCRMSocialConfig = `{
  "chart":{"height":300,"type":"bar","background":"transparent","toolbar":{"show":false}},
  "plotOptions":{"bar":{"horizontal":true,"distributed":true,"borderRadius":2,"borderRadiusApplication":"end"}},
  "stroke":{"colors":["var(--color-base-100)"]},
  "colors":["#5860ff"],
  "series":[{"data":[48,35,10,7]}],
  "xaxis":{"categories":["Facebook","Whatsapp","Instagram","Youtube"],"type":"category"},
  "grid":{"show":false}
}`

const nxGlobalSalesConfig = `{
  "chart":{"height":344,"type":"bar","parentHeightOffset":0,"background":"transparent","toolbar":{"show":false}},
  "plotOptions":{"bar":{"horizontal":true,"borderRadius":4,"distributed":true,"borderRadiusApplication":"end"}},
  "dataLabels":{"enabled":true,"textAnchor":"start","style":{"colors":["#fff"]},"offsetX":-10,"dropShadow":{"enabled":false}},
  "series":[{"data":[9,12,13,16,14,17,19]}],
  "legend":{"show":false},
  "stroke":{"width":0,"colors":["#fff"]},
  "xaxis":{"categories":["Turkey","India","Canada","US","Netherlands","Italy","Other"]},
  "yaxis":{"labels":{"show":false}},
  "grid":{"show":false},
  "tooltip":{"theme":"dark","x":{"show":false}},
  "colors":["#7179ff","#4bcd89","#ff6c88","#5cb7ff","#9071ff","#ff5892","#ff8b4b"]
}`
