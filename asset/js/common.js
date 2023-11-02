$(document).ready(function () {
	//side-menu
	$(".subMenuButton").on("click", function (event) {
		event.preventDefault();
		var $target = $(this).closest('li').children('ul');

		$target.slideToggle();
		$(this).toggleClass("fa-plus fa-minus");
	});

	$(".leftMenuButton").on("click", function (event) {
		var $target = $(".leftMenu");
		$target.toggleClass("leftMenuShow");
		$(this).toggleClass("leftMenuButtonMove");
	});

	$("main").on("click", ".buttonClose", function () {
		var $alert = $(this).parents(".alert");

		$alert.fadeOut(function () {
			$alert.fadeOut(function () {
				$alert.remove();
				var $alertList = $(".alert");
				if ($alertList.length == 0) {
					$(".siteMessage").remove();
				}
			});
		});
	});

	setInterval(function () {
		var $autoClear = $(".alert.autoClear");

		if ($autoClear.length != 0) {

			$autoClear.fadeOut(function () {
				$autoClear.remove();
				var $alertList = $(".alert");
				if ($alertList.length == 0) {
					$(".siteMessage").remove();
				}
			});
		}
	}, 8000);

	//pagination
	$("ul.pagination li.disabled").on("click", function () {
		return false;
	});

	//header
	var $leftMenu = $(".leftMenu");
	var $leftMenuButton = $(".leftMenuButton");

	var $topMenuButton = $(".topMenuButton");
	var $topMenu = $(".topMenu");
	var $formSearch = $(".formSearch");
	var $header = $("header");
	var $main = $("main");

	if ($header.length) {
		var flex = new $.FlexWatch({ 'debug': false, 'update': updateHeight });

		if ($formSearch.length) {
			$formSearch.on("dropped", function (event) {
				$topMenuButton.show();
				$formSearch.css("margin-top", "1.2rem");
			});

			$formSearch.on("risen", function () {
				$topMenuButton.hide();
				$formSearch.css("margin-top", "0");
			});

			flex.addWatch(".formSearch", ".siteTitle");
		}

		if ($topMenu.length) {
			$topMenu.on("dropped", function () {
				$topMenuButton.show();
				$topMenu.css("margin-top", "1.2rem");
			});

			$topMenu.on("risen", function () {
				$topMenuButton.hide();
				$topMenu.css("margin-top", "0");
			});

			flex.addWatch(".topMenu", ".siteTitle");
		}

		flex.startWatching();
	}

	$topMenuButton.on("click", function () {
		$header.toggleClass("autoHeight");
		updateHeight();
	});

	function updateHeight() {
		if ($header.length == 0) {
			return;
		}

		var headerHeight = $header.outerHeight();
		$topMenuButton.css("top", headerHeight + "px");
		$leftMenuButton.css("top", headerHeight + "px");
		$leftMenu.css("top", headerHeight + "px");
		$main.css("padding-top", headerHeight + "px");
	}

	//radio button labels
	$("input[type='radio']").siblings("span.radioLabel").on("click", function () {
		var $radio = $(this).parent().children("input[type='radio']");
		$radio.prop("checked", !$radio.prop("checked"));
		$radio.trigger("change");
	});

	$("input[type='checkbox']").siblings("span.radioLabel").on("click", function () {
		var $radio = $(this).parent().children("input[type='checkbox']");
		$radio.prop("checked", !$radio.prop("checked"));
		$radio.trigger("change");
	});

	//product small images
	$("img.smallImage").on("click", function () {
		$(this).parent().find("a.popupImage").modal("open");
		return false;
	});
});

function showAjaxMsg(message) {
	var $messageDiv = $(".main div.siteMessage");

	var str = "<div class=\"row\">";
	str += "<div class=\"col\">";
	str += message;
	str += "</div>";
	str += "</div>";
	$("main div.siteMessage").html(str);
}