$(document).ready(function() {
  GetBudgets();
});

var GetBudgets = function() {
  const data = {action: "getbudgets"};
  request(data, (res)=>{
    if (!!res && !!res.message && res.message.length > 0) {
      budgetList = JSON.parse(res.message);
    }
    budgetList.forEach(v => {
      $("#budgets").append("<a href='#' onclick='GetBudget(\"" + v + "\");return false;'>" + v + "</a><span>&nbsp;</span>");
    });
  }, (e)=>{
    console.log(e.responseJSON.message);
    $("#warning").text(e.responseJSON.message).removeClass("hidden").addClass("visible");
  });
};

var GetBudget = function(budgetName) {
  const data = {action: "getbudget", name: budgetName};
  request(data, (res)=>{
    if (!!res && !!res.message && res.message.length > 0) {
      budgetData = JSON.parse(res.message);
    }
    $("#result").text(JSON.stringify(budgetData, null, "\t"));
    $("#info").removeClass("hidden").addClass("visible");
  }, (e)=>{
    console.log(e.responseJSON.message);
    $("#warning").text(e.responseJSON.message).removeClass("hidden").addClass("visible");
  });
};

var request = function(data, callback, onerror) {
  $.ajax({
    type:          'POST',
    dataType:      'json',
    contentType:   'application/json',
    scriptCharset: 'utf-8',
    data:          JSON.stringify(data),
    url:           App.url
  })
  .done(function(res) {
    callback(res);
  })
  .fail(function(e) {
    onerror(e);
  });
};
var App = { url: location.origin + {{ .ApiPath }} };
