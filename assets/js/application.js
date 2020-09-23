require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
require("@fortawesome/fontawesome-free/js/all.js");

$(() => {
    $(".complete").on("click", (event) => {
        $.ajax({
            method: "POST",
            url: "http://habitbeds.com/beds/toggle_complete",
            data: { bedid: event.target.dataset.bedid },
        });
    });
});
