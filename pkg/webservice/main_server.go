/*
Copyright (C) 2018 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package webservice

import (
	"strings"

	"github.com/blackducksoftware/perceptor-protoform/pkg/api/hub/v1"
	"github.com/blackducksoftware/perceptor-protoform/pkg/hub"
	model "github.com/blackducksoftware/perceptor-protoform/pkg/model"
	"github.com/gin-gonic/contrib/static"
	gin "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SetupHTTPServer will used to create all the http api
func SetupHTTPServer(hc *hub.Creater, config *model.Config) {
	go func() {
		// data, err := ioutil.ReadFile("/public/index.html")
		// Set the router as the default one shipped with Gin
		router := gin.Default()
		// Serve frontend static files
		router.Use(static.Serve("/", static.LocalFile("/views", true)))

		// prints debug stuff out.
		router.Use(GinRequestLogger())

		router.GET("/api/sql-instances", func(c *gin.Context) {
			keys := []string{"pvc-000", "pvc-001", "pvc-002"}
			c.JSON(200, keys)
		})

		router.GET("/hub", func(c *gin.Context) {
			log.Debug("get hub request")
			hubs, err := hc.HubClient.SynopsysV1().Hubs(corev1.NamespaceDefault).List(metav1.ListOptions{})
			if err != nil {
				log.Errorf("unable to get the hub list due to %+v", err)
				c.JSON(404, "\"message\": \"Failed to List the hub\"")
			}

			returnVal := make(map[string]*v1.Hub)

			for _, v := range hubs.Items {
				returnVal[v.Spec.Namespace] = &v
			}

			c.JSON(200, returnVal)
		})

		router.POST("/hub", func(c *gin.Context) {
			log.Debug("create hub request")
			hubSpec := &v1.HubSpec{}
			if err := c.BindJSON(hubSpec); err != nil {
				log.Debugf("Fatal failure binding the incoming request ! %v", c.Request)
			}
			hubSpec.State = "pending"

			if strings.EqualFold(hubSpec.PostgresPassword, "") {
				hubSpec.PostgresPassword = "blackduck"
			}

			ns, err := hc.KubeClient.CoreV1().Namespaces().Create(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Namespace: hubSpec.Namespace, Name: hubSpec.Namespace}})
			log.Debugf("created namespace: %+v", ns)
			if err != nil {
				log.Errorf("unable to create the namespace due to %+v", err)
				c.JSON(404, "\"message\": \"Failed to create the namespace\"")
			}
			hc.HubClient.SynopsysV1().Hubs(hubSpec.Namespace).Create(&v1.Hub{ObjectMeta: metav1.ObjectMeta{Name: hubSpec.Namespace}, Spec: *hubSpec})

			c.JSON(200, "\"message\": \"Succeeded\"")
		})

		router.DELETE("/hub", func(c *gin.Context) {
			var hubSpec string
			if err := c.BindJSON(hubSpec); err != nil {
				log.Debugf("Fatal failure binding the incoming request ! %v", c.Request)
			}

			log.Debugf("delete hub request %v", hubSpec)

			// This is on the event loop.
			hc.HubClient.SynopsysV1().Hubs(hubSpec).Delete(hubSpec, &metav1.DeleteOptions{})

			c.JSON(200, "\"message\": \"Succeeded\"")
		})

		// Start and run the server - blocking call, obviously :)
		router.Run(":80")
	}()
}
