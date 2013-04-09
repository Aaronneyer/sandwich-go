$(document).ready(function(){
    $("#query_form").on("submit", function(e) {
        e.preventDefault();
        $.get("/search", {search:
            $(this).find("input[type=text]").val()}, function(data){
                $("#content").html(data);
                $("#content").trigger("change");
            });
    })
    $("#killbtn").on("click", function(e) {
      e.preventDefault();
      $.get("/kill")
      return false;
    });
    $("#updatebtn").on("click", function(e) {
      e.preventDefault();
      $.get("/update")
      return false;
    });
    x = function(){ $(".dl-link").on("click", function(e) {
        e.preventDefault();
        $.get("/download", {ip: $(this).attr("data-ip"), file: $(this).attr("data-file")});
    })}
    $("#content").on("change", x);
})
