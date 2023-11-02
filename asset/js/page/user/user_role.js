$(document).ready(function () {
	var timeout;

	$(".revokeButton").on("click", function () {
		$(this).parent().hide();
		$(this).parent().next(".revokeConfirm").show();

		var self = this;
		timeout = setTimeout(function () {
			$(self).parent().show();
			$(self).parent().parent().find(".revokeConfirm").hide();
		}, 10000);
		return false;
	});

	$(".cancelButton").on("click", function () {
		clearTimeout(timeout);

		$(this).parent().hide();
		$(this).parent().parent().find(".revoke").show();
	});
})