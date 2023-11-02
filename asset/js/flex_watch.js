(function ($) {

	$.FlexWatch = function (opts) {
		if (opts === undefined) {
			opts = {}
		}

		var defaults = {
			interval: 50,
			debug: false
		}

		this.options = $.extend(defaults, opts)
		this.watchList = [];

		this.options.update();
	};

	$.FlexWatch.prototype.addWatch = function (target, reference) {
		var $tar = $(target);
		if ($tar.length === 0) {
			console.log("could not find target:", $tar.selector);
			return;
		}

		var $ref = $(reference);
		if ($ref.length === 0) {
			console.log("could not find reference:", $ref.selector);
			return;
		}

		this.watchList.push({
			$tar: $tar,
			$ref: $ref
		});
	};

	$.FlexWatch.prototype.getState = function ($tar, $ref) {
		var tarTop = $tar.offset().top;
		var refTop = $ref.offset().top;
		var refHeight = $ref.outerHeight(true);

		if (tarTop >= refTop + refHeight) {
			return "dropped";
		} else {
			return "risen";
		}
	};

	$.FlexWatch.prototype.checkState = function () {
		var watchListLen = this.watchList.length
		for (var i = 0; i < watchListLen; i++) {
			var state = this.getState(this.watchList[i].$tar, this.watchList[i].$ref);

			if (state === this.watchList[i].state) {
				continue;
			}

			if (this.options.debug) {
				console.log("triggering for ", this.watchList[i].$tar.selector, state);
			}

			this.watchList[i].$tar.trigger(state);
			this.watchList[i].state = state;
		}
	};

	$.FlexWatch.prototype.startWatching = function () {
		if (this.started) {
			console.log("flesWatch has already been started.");
			return;
		}

		this.checkState();

		var self = this;
		$(window).on("resize", function () {
			setTimeout(function () {
				self.checkState();
				self.options.update();
			}, self.options.interval);
		});
	};
})(jQuery);