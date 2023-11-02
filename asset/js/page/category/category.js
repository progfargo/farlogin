$(document).ready(function () {
	var $allParent = $("td[class*='parent-']");
	var $allSub = $("td[class*='sub-']");

	$allParent.on("click", function () {
		var classStr = $(this).attr("class");

		var reg = /parent-(\d+)/g;
		var match = reg.exec(classStr);
		var parentId = match[1];

		$collapseIcon = $(this).children("span.collapseIcon");
		$expandIcon = $(this).children("span.expandIcon");

		$targetSub = $("td.sub-" + parentId).parent();
		if ($collapseIcon.is(":visible")) {
			$collapseIcon.hide();
			$expandIcon.show();
			$targetSub.find("span.collapseIcon").hide();
			$targetSub.find("span.expandIcon").show();
			$targetSub.hide();
		} else {
			$collapseIcon.show();
			$expandIcon.hide();
			$targetSub.find("span.collapseIcon").show();
			$targetSub.find("span.expandIcon").hide();
			$targetSub.show();
		}
	});

	$(".expandAll").on("click", function () {
		expandAll();
	});

	function expandAll() {
		$allSub.parent().not(".sub-root").show();
		$allParent.children("span.collapseIcon").show();
		$allParent.children("span.expandIcon").hide();
	}

	$(".collapseAll").on("click", function () {
		collapseAll();
	});

	function collapseAll() {
		$allSub.parent().not(".sub-root").hide();
		$allParent.children("span.collapseIcon").hide();
		$allParent.children("span.expandIcon").show();
	}

	collapseAll();

	$(".buttonDelete").on("click", function () {
		$(this).parent().hide();
		$(this).parent().next(".deleteConfirm").show();
		return false;
	});

	$(".cancelButton").on("click", function () {
		$(this).parent().hide();
		$(this).parent().parent().find(".editDelete").show();
	});
})