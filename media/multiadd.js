(function() {
    this.addFactForm = function(idx) {
        var el = $("#fact-row-0").clone();
        el.attr('id', 'fact-row-' + idx);
        $("#fact-row-" + (idx - 1) + " .title").unbind('focus');
        $("#insert-point").before(el);

        $("#fact-row-" + idx + " .title")
          .attr('name', 'title' + idx).val('');
        $("#fact-row-" + idx + " .details")
          .attr('name', 'details' + idx).val('');
        $("#fact-row-" + idx + " .source-name").attr({
            'name': 'source_name' + idx,
            'value': $("#fact-row-" + (idx - 1) + " .source-name").val()
        });
        $("#fact-row-" + idx + " .source-url").attr({
            'name': 'source_url' + idx,
            'value': $("#fact-row-" + (idx - 1) + " .source-url").val()
        });
        $("#fact-row-" + idx + " .tags").attr({
            'name': 'tags' + idx,
            'value': $("#fact-row-" + (idx - 1) + " .tags").val()
        });

        $("#fact-row-" + idx + " .title").bind('focus', function(evt) {
            if(evt.target.value === "") {
                window.addFactForm(idx + 1);
            }});
    };

	 $("#fact-row-0 .title").bind('focus', function(evt) {
		 if(evt.target.value === "") {
			 window.addFactForm(1);
		 }
		});
 }());
