$(document).ready(function () {
	$("select[name='catId']").change(function () {
		$("#categoryForm").submit();
	});
});