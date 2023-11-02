$(document).ready(function () {
	var $contactTypeForm = $("#contactTypeForm");

	$("select[name='contactType']").on("change", function () {
		console.log("test");
		$contactTypeForm.submit();

		return false;
	});
});