import React, { useRef, useEffect, useState } from "react";
import * as d3 from "d3";

function Graph() {
    const svgRef = useRef();
    useEffect(() => {
        var svg = d3.select(svgRef.current);
        svg.selectChildren().remove();
        svg = svg.attr("width", 200)
            .attr("height", 200)
            .style("padding", 2)
            .style("margin", 2)
            .append("g")
            .attr("transform", "translate(25,25)");
        if (this.props.data.length == 0) {
            return;
        }
        var h = 150;
        var w = 150;
        var xScale = d3.scaleLinear()
            .domain([0, this.props.data.length - 1])
            .rangeRound([0, w]);
        svg.append("g")
            .attr("class", "axis")
            .attr("transform", "translate(0," + h + ")")
            .call(d3.axisBottom()
                .ticks(4)
                .scale(xScale));

        var yScale = d3.scaleLinear()
            .domain(d3.extent(this.props.data))
            .rangeRound([h, 0]);
        svg.append("g")
            .attr("class", "axis")
            .call(d3.axisLeft()
                .ticks(4)
                .scale(yScale));

        svg.append("path")
            .datum(this.props.data)
            .attr("fill", "none")
            .attr("stroke-width", 1.5)
            .attr("d", d3.line().x((d, i) => xScale(i)).y((d, i) => yScale(d)))
            .attr("stroke", "black");
    }, [this.props.data]);
    return (
        <React.Fragment>
            <svg ref={svgRef}></svg>
        </React.Fragment>
    );
}

export default Graph;