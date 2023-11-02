$(document).ready(function () {
	$(".copyLink").on("click", function (event) {
		var url = new URL(location.href);
		var linkStr = url.origin + $(this).prev("a").attr("href");

		var $textArea = $("<textarea class=\"copyClipboard\">" + linkStr + "</textarea>");
		$textArea.insertAfter("body");
		$textArea.select();
		document.execCommand("copy");
		$(this).hide().fadeIn();
		$(".copyClipboard").remove();
	});
});