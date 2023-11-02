$(document).ready(function() {
	$(".copyLink").on("click", function(event) {
		var str = $(this).prev("a").attr("href");
		
		var $textArea = $("<textarea class=\"copyClipboard\">" + str + "</textarea>");
		$textArea.insertAfter("body");
		$textArea.select();
		document.execCommand("copy");
		$(this).hide().fadeIn();
		$(".copyClipboard").remove();
	});
});