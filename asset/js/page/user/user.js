$(document).ready(function () {
	$("select[name='rid']").change(function () {
		$("#ridForm").submit();
	});

	$("select[name='stat']").change(function () {
		$("#statForm").submit();
	});

	$("select[name='domain']").change(function () {
		$("#domainForm").submit();
	});
});
