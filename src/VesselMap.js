import React, { useRef, useEffect } from "react";
import ArcGISMap from "@arcgis/core/Map";
import axios from 'axios';
import BasemapToggle from "@arcgis/core/widgets/BasemapToggle";
import Chart from 'chart.js/auto';
import esriConfig from '@arcgis/core/config.js';
import FeatureLayer from "@arcgis/core/layers/FeatureLayer";
import LayerList from "@arcgis/core/widgets/LayerList";
import Legend from "@arcgis/core/widgets/Legend";
import MapView from "@arcgis/core/views/MapView";
import "./VesselMap.css";

function VesselMap(props){

  esriConfig.assetsPath = './assets';
  const mapDiv = useRef(null);

  useEffect(() => {
    if (mapDiv.current){
      const popActionTS = {
        title: "Time Series",
        id: "timeSeries"
      }
      const popActionUV = {
        title: "Vessel Type by Unique Vessel",
        id: "unqVessel"
      }
      const popActionTN = {
        title: "Vessel Type by trip number",
        id: "tripNum"
      } 
      // Template for popup for structure points
      const structureTemplate = {
        title: "{Name}",
        content: [{
          type: "fields",
          fieldInfos: [
            {
              fieldName: "Year",
              label: "Year Constructed"
            },
            {
              fieldName: "Type",
              label: "Structure Type"
            },
            {
              fieldName: "Length",
              label: "Structure Length"
            },
            {
              fieldName: "Community",
              label: "Community Number"
            },
            {
              fieldName: "Count",
              label: "Gross Traffic Count"
            }
          ]
        }],
        actions: [popActionTS,popActionUV,popActionTN]
      }
  
      // Tells the structure layer how to render the points
      const youngStructure = new Date().getFullYear() - 2000;
      const middleStructure = new Date().getFullYear() - 1971;
      const oldStructure = new Date().getFullYear() - 1970;
      const pointRenderer = {
        type: "class-breaks",
        field: "Community",
        defaultSymbol: {
          type: "simple-marker",
          color: "white",
          outline: {
            width: 0.5,
            color: "white"
          },
          size: 10,
        },
        classBreakInfos: [
          {
            minValue: 0,
            maxValue: 0,
            symbol: {
              type: "simple-marker",
              color: "white",
              style: "circle",
              outline: {
                width: 0.5,
                color: "white"
              },
              size: 10,
            }
          },
          {
            minValue: 1,
            maxValue: 1,
            symbol: {
              type: "simple-marker",
              color: "white",
              style: "cross",
              outline: {
                width: 0.5,
                color: "white"
              },
              size: 10,
            }
          },
          {
            minValue: 2,
            maxValue: 2,
            symbol: {
              type: "simple-marker",
              color: "white",
              style: "diamond",
              outline: {
                width: 0.5,
                color: "white"
              },
              size: 10,
            }
          },
          {
            minValue: 3,
            maxValue: 3,
            symbol: {
              type: "simple-marker",
              color: "white",
              style: "square",
              outline: {
                width: 0.5,
                color: "white"
              },
              size: 10,
            }
          },
          {
            minValue: 4,
            maxValue: 4,
            symbol: {
              type: "simple-marker",
              color: "white",
              style: "triangle",
              outline: {
                width: 0.5,
                color: "white"
              },
              size: 10,
            }
          },
        ],
        visualVariables: [
          {
            type: "color",
            valueExpressionTitle: "Age in years",
            valueExpression: "2021-$feature.Year",
            stops: [
              { value: oldStructure, color: "red" },
              { value: middleStructure, color: "yellow" },
              { value: youngStructure, color: "green" },
            ]
          },
          {
            type: "opacity",
            field: "Count",
            stops: [
              { value: 0, opacity: 0.4 },
              { value: 5000, opacity: 1 },
            ]
          }
        ]
      }
  
      const structureLayer = new FeatureLayer({
        title: "Structures",
        source: props.points, //Tell the layer where to get the data for the points
        renderer: pointRenderer,
        popupTemplate: structureTemplate,
        objectIDField: "ObjectID",
        fields: [
          {
            name: "ObjectID",
            type: "oid"
          },
          {
            name: "Name",
            type: "string"
          },
          {
            name: "Year",
            type: "integer"
          },
          {
            name: "Type",
            type: "string"
          },
          {
            name: "Length",
            type: "integer"
          },
          {
            name: "Community",
            type: "integer"
          },
          {
            name: "Count",
            type: "integer"
          }
        ]
      });

      const map = new ArcGISMap({
        basemap: "topo-vector",  // initial map styling
        layers: [props.divisionLayer, props.districtLayer, structureLayer]  // array of layers that sits on top of the basemap
      });

      const view = new MapView({
        map: map,
        container: mapDiv.current,
        // sets the initial positioning and zoom of the map
        extent: {
          spatialReference: {
            wkid: 102100,
          },
          xmax: -47458864.13281338,
          xmin: -54033671.557789356,
          ymax: 6678183.73641666,
          ymin: 1937864.9902844299,
        },
        zoom: 4
      });

      function timeSeries(e) {
        var id = view.popup.selectedFeature.attributes.ObjectID;
        axios.get(`/api/timeseries/${id}`)
          .then(function (response) {
            view.popup.visible = true;
            view.popup.open({
              title: "Time Series Graph",
              content: setContentInfoTimeSeries(response.data)
            });
          });
      }
      function setContentInfoTimeSeries(response) {
        var canvas = document.createElement('canvas');
        canvas.id = "timeSeries";
        var sum = [];
        var days = [];
        for (var key in response) {
          sum.push(response[key]['Sum']);
          days.push(response[key]['Day']);
        }
        var data = {
          datasets: [{
            label: "Vessel Trip Count",
            data: sum,
            backgroundColor: ["#4286f4", "#41f4be", "#8b41f4", "#e241f4", "#f44185", "#f4cd41"]
          }],
          labels: days
        };
        var timeSeriesChart = new Chart(canvas, {
          type: 'bar',
          data: data,
          options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
              y: {
                beginAtZero: true
              }
            }
          },
        });
        return canvas;
      }
      function unqVessel(e) {
        var id = view.popup.selectedFeature.attributes.ObjectID;
        axios.get(`/api/uniquevessels/${id}`)
          .then(function (response) {
            view.popup.visible = true;
            view.popup.open({
              title: "Distribution of Vessel Type by unique vessels",
              content: setContentInfoUnqVessel(response.data)
            });
          });
      }
      function setContentInfoUnqVessel(response) {
        var canvas = document.createElement('canvas');
        canvas.id = "unqVessel";
        var count = [];
        var vessels = [];
        for (var key in response) {
          count.push(response[key]['Count']);
          vessels.push(response[key]['Vessel']);
        }
        var data = {
          datasets: [{
            label: "Vessel Type",
            data: count,
            backgroundColor: ["#4286f4", "#41f4be", "#8b41f4", "#e241f4", "#f44185", "#f4cd41"]
          }],
          labels: vessels
        };
        var unqVesselChart = new Chart(canvas, {
          type: 'bar',
          data: data,
          options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
              y: {
                beginAtZero: true
              }
            }
          },
        });
        return canvas;
      }
      function tripNum(e) {
        var id = view.popup.selectedFeature.attributes.ObjectID;
        axios.get(`/api/vesseltripcounts/${id}`)
          .then(function (response) {
            view.popup.visible = true;
            view.popup.open({
              title: "Distribution of Vessel Type by number of trips",
              content: setContentInfoTripCount(response.data)
            });
          });
      }
      function setContentInfoTripCount(response) {
        var canvas = document.createElement('canvas');
        canvas.id = "tripNum";
        var count = [];
        var vessels = [];
        for (var key in response) {
          count.push(response[key]['Count']);
          vessels.push(response[key]['Vessel']);
        }
        var data = {
          datasets: [{
            label: "Vessel Type",
            data: count,
            backgroundColor: ["#4286f4", "#41f4be", "#8b41f4", "#e241f4", "#f44185", "#f4cd41"]
          }],
          labels: vessels
        };
        var tripCountChart = new Chart(canvas, {
          type: 'bar',
          data: data,
          options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
              y: {
                beginAtZero: true
              }
            }
          },
        });
        return canvas;
      }
      view.popup.on("trigger-action", function (event) {
        if (event.action.id === "timeSeries") {
          timeSeries(event);
        }
        else if (event.action.id === "unqVessel") {
          unqVessel(event);
        }
        else if (event.action.id === "tripNum") {
          tripNum(event);
        }
      });
      view.when(function () {
        var layerList = new LayerList({
          view: view,
        });
        // adds the Layer toggle to the top-right of the screen
        view.ui.add(layerList, "top-right");
      });

      //BasemapToggle object is used to toggle the "style" of the basemap
      var toggle = new BasemapToggle({
        view: view, // view that provides access to the map's 'topo-vector' basemap
        nextBasemap: "hybrid" // allows for toggling to the 'hybrid' basemap
      });
      // adds the toggle base map button display
      view.ui.add(toggle, "top-right");

      view.ui.add(new Legend({
        view: view
      }), "bottom-left");

    }
  },[props.divisionLayer, props.districtLayer, props.points]);


  return (<div className="mapDiv" ref={mapDiv}></div>);
}

export default VesselMap;