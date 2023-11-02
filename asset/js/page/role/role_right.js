$(document).ready(function () {
	$(".selectAll").on("click", function (event) {
		event.stopPropagation();
		$(this).closest("table").find(":checkbox").prop("checked", true);
	});

	$(".selectNone").on("click", function (event) {
		event.stopPropagation();
		$(this).closest("table").find(":checkbox").prop("checked", false);
	});

	$(".selectReverse").on("click", function (event) {
		event.stopPropagation();
		var list = $(this).closest("table").find(":checkbox");
		list.each(function (i) {
			this.checked = !this.checked
		});
	});

	var $tableBody = $("table tbody");
	var $selectAll = $(".selectAll");
	var $selectNone = $(".selectNone");
	var $selectReverse = $(".selectReverse");

	var $expandAll = $(".expandAll");
	var $collapseAll = $(".collapseAll");

	$collapseAll.on("click", function () {
		$tableBody.fadeOut("fast");
		$selectAll.hide();
		$selectNone.hide();
		$selectReverse.hide();
	});

	$expandAll.on("click", function () {
		$tableBody.fadeIn("fast");
		$selectAll.show();
		$selectNone.show();
		$selectReverse.show();
	});

	$("tr.roleHeader").on("click", function () {
		$(this).closest("table").find("tbody").fadeToggle("fast");
		$(this).find(".selectAll").toggle();
		$(this).find(".selectNone").toggle();

		var $selectReverse = $(this).find(".selectReverse");
		$selectReverse.toggle();

		if ($selectReverse.is(":visible")) {
			$(this).find(".accordionIcon").removeClass("fa-caret-down").addClass("fa-caret-up");
		} else {
			$(this).find(".accordionIcon").removeClass("fa-caret-up").addClass("fa-caret-down");
		}
	});
});
