$(document).ready(function () {
	var timeout;

	$(".deleteButton").on("click", function () {
		console.log("test");
		$(this).parent().hide();
		$(this).parent().next(".deleteConfirm").show();

		var self = this;
		timeout = setTimeout(function () {
			$(self).parent().show();
			$(self).parent().parent().find(".deleteConfirm").hide();
		}, 10000);
		return false;
	});

	$(".cancelButton").on("click", function () {
		clearTimeout(timeout);

		$(this).parent().hide();
		$(this).parent().parent().find(".delete").show();
	});

	$(".copyLink").on("click", function (event) {
		var str = $(this).parents("td").find("span.sessionHash").text();

		var $textArea = $("<textarea class=\"copyClipboard\">" + str + "</textarea>");
		$textArea.insertAfter("body");
		$textArea.select();
		document.execCommand("copy");
		$(this).hide().fadeIn();
		$(".copyClipboard").remove();
	});
})