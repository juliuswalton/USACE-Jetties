import React, { useRef, useEffect } from "react";
//import FeatureLayer from "@arcgis/core/layers/FeatureLayer";
import ArcGISMap from "@arcgis/core/Map";
//import DictionaryRenderer from "@arcgis/core/renderers/DictionaryRenderer";
import MapView from "@arcgis/core/views/MapView";
import GeoJSONLayer from "@arcgis/core/layers/GeoJSONLayer";
import LayerList from "@arcgis/core/widgets/LayerList";
import BasemapToggle from "@arcgis/core/widgets/BasemapToggle";
import esriConfig from '@arcgis/core/config.js';
import axios from 'axios';
import FeatureLayer from "@arcgis/core/layers/FeatureLayer";
import Legend from "@arcgis/core/widgets/Legend";
import Chart from 'chart.js/auto';

import "./App.css";
axios.defaults.baseURL = "http://localhost:8080"; //<== USE THIS LINE FOR DEVELOPMENT ON LOCAL MACHINE
//axios.defaults.baseURL = "https://strange-tome-305601.ue.r.appspot.com/"; //<== USE THIS LINE FOR PRODUCTION
function App() {

  // Required: Set this property to insure assets resolve correctly.
  esriConfig.assetsPath = './assets';

  const mapDiv = useRef(null);

  useEffect(() => {
    if (mapDiv.current) {
      /**
       * Initialize application
       */

      // Template for the Division's popup
      const divisionTemplate = {
        title: "USACE Civil Works Division",
        content: "{DIV_SYM} : {DIVISION}",
      };
      // Layer that pulls in the geoJSON data for the USACE Divisions
      const divisionLayer = new GeoJSONLayer({
        url: "https://opendata.arcgis.com/datasets/4cdab8820d7a4f58aaa003be63f059ac_0.geojson",
        copyright: "USACE Civil Works",
        popupTemplate: divisionTemplate,
        title: "USACE Divisions",
        listMode: false,
        visible: false,
        legendEnabled: false
      });

      const districtTemplate = {
        title: "USACE Civil Works District",
        content: "{DIV_SYM} : {District}",
      };
      const districLayer = new GeoJSONLayer({
        url: "https://opendata.arcgis.com/datasets/f3e0ba4566094e74910c391eb4ecc99f_0.geojson",
        copyright: "USACE Civil Works",
        popupTemplate: districtTemplate,
        title: "USACE Districts",
        listMode: false,
        visible: false,
        legendEnabled: false
      });


      // Creates the map component
      const map = new ArcGISMap({
        basemap: "topo-vector",  // initial map styling
        layers: [divisionLayer, districLayer]       // array of layers that sits on top of the basemap
      });

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

      const youngStructure = new Date().getFullYear() - 20;
      const middleStructure = new Date().getFullYear() - 50;
      const oldStructure = new Date().getFullYear() - 51;
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
            field: "Year",
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
              { value: 0, opacity: 0.2 },
              { value: 5000, opacity: 1 },
            ]
          }
        ]
      }




      //Api call to get the time series data for total counts of transits for each day
      /* axios.get('/api/timeseries')
       .then(function(response){
         console.log(response); });
 
 
       axios.get('/api/uniquevessels')
       .then(function(response){
         console.log(response); });
 
       axios.get('/api/vesseltripcounts')
       .then(function(response){
         console.log(response); });
 */

      // Api call to get structure data
      var structurePoints = [];
      axios.get('/api/structure')
        .then(function (response) {
          console.log(response);
          var structures = response.data; //Grab response data

          //For each point in the response data create a ArcGIS Point Graphic
          for (var i = 0; i < structures.length; i++) {
            var feature = {
              geometry: {
                type: "point",
                x: structures[i].Lon,
                y: structures[i].Lat
              },
              attributes: {
                ObjectID: structures[i].ID,
                Name: structures[i].Name,
                Year: structures[i].Year,
                Type: structures[i].Type,
                Length: structures[i].Length,
                Community: structures[i].Community,
                Count: structures[i].Count
              }
            };

            //Add Point to the array of points
            structurePoints.push(feature);
          }

          //Create layer to show the structures
          const structureLayer = new FeatureLayer({
            title: "Structures",
            source: structurePoints, //Tell the layer where to get the data for the points
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

          map.add(structureLayer); //Add the layer to the base map
        });
      // Creates the MapView - Necessary for rendering the Map Object above
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
        console.log(view.popup.selectedFeature)
        var id = view.popup.selectedFeature.attributes.ObjectID;
        axios.get(`/api/timeseries/${id}`)
          .then(function (response) {
            console.log("TIME SERIES", response);
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
        console.log(sum);
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
        console.log(view.popup.selectedFeature)
        var id = view.popup.selectedFeature.attributes.ObjectID;
        axios.get(`/api/uniquevessels/${id}`)
          .then(function (response) {
            console.log("UNIQUE VESSEL", response);
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
        console.log(view.popup.selectedFeature)
        var id = view.popup.selectedFeature.attributes.ObjectID;
        axios.get(`/api/vesseltripcounts/${id}`)
          .then(function (response) {
            console.log("TRIP COUNT", response);
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
      //Creates the LayerList which is used to toggle visiblity of available map layers
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
  }, []);

  return (<div className="mapDiv" ref={mapDiv}></div>

  );
}

export default App;
