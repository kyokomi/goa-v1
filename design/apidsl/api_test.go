package apidsl_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kyokomi/goa-v1/design"
	"github.com/kyokomi/goa-v1/design/apidsl"
	"github.com/kyokomi/goa-v1/dslengine"
)

var _ = Describe("API", func() {
	var name string
	var dsl func()

	BeforeEach(func() {
		dslengine.Reset()
		name = ""
		dsl = nil
	})

	JustBeforeEach(func() {
		apidsl.API(name, dsl)
		dslengine.Run()
	})

	Context("with no DSL", func() {
		BeforeEach(func() {
			name = "foo"
		})

		It("produces a valid API definition", func() {
			Ω(Design.Validate()).ShouldNot(HaveOccurred())
			Ω(Design.Name).Should(Equal(name))
		})
	})

	Context("with an already defined API with the same name", func() {
		BeforeEach(func() {
			name = "foo"
		})

		It("produces an error", func() {
			apidsl.API(name, dsl)
			Ω(dslengine.Errors).Should(HaveOccurred())
		})
	})

	Context("with an already defined API with a different name", func() {
		BeforeEach(func() {
			name = "foo"
		})

		It("returns an error", func() {
			apidsl.API("news", dsl)
			Ω(dslengine.Errors).Should(HaveOccurred())
		})
	})

	Context("with valid DSL", func() {
		JustBeforeEach(func() {
			Ω(dslengine.Errors).ShouldNot(HaveOccurred())
			Ω(Design.Validate()).ShouldNot(HaveOccurred())
		})

		Context("with a description", func() {
			const description = "description"

			BeforeEach(func() {
				dsl = func() {
					apidsl.Description(description)
				}
			})

			It("sets the API description", func() {
				Ω(Design.Description).Should(Equal(description))
			})
		})

		Context("with a title", func() {
			const title = "title"

			BeforeEach(func() {
				dsl = func() {
					apidsl.Title(title)
				}
			})

			It("sets the API title", func() {
				Ω(Design.Title).Should(Equal(title))
			})
		})

		Context("with a version", func() {
			const version = "2.0"

			BeforeEach(func() {
				dsl = func() {
					apidsl.Version(version)
				}
			})

			It("sets the API version", func() {
				Ω(Design.Version).Should(Equal(version))
			})
		})

		Context("with a terms of service", func() {
			const terms = "terms"

			BeforeEach(func() {
				dsl = func() {
					apidsl.TermsOfService(terms)
				}
			})

			It("sets the API terms of service", func() {
				Ω(Design.TermsOfService).Should(Equal(terms))
			})
		})

		Context("with contact information", func() {
			const contactName = "contactName"
			const contactEmail = "contactEmail"
			const contactURL = "http://contactURL.com"

			BeforeEach(func() {
				dsl = func() {
					apidsl.Contact(func() {
						apidsl.Name(contactName)
						apidsl.Email(contactEmail)
						apidsl.URL(contactURL)
					})
				}
			})

			It("sets the contact information", func() {
				Ω(Design.Contact).Should(Equal(&ContactDefinition{
					Name:  contactName,
					Email: contactEmail,
					URL:   contactURL,
				}))
			})
		})

		Context("with license information", func() {
			const licenseName = "licenseName"
			const licenseURL = "http://licenseURL.com"

			BeforeEach(func() {
				dsl = func() {
					apidsl.License(func() {
						apidsl.Name(licenseName)
						apidsl.URL(licenseURL)
					})
				}
			})

			It("sets the API license information", func() {
				Ω(Design.License).Should(Equal(&LicenseDefinition{
					Name: licenseName,
					URL:  licenseURL,
				}))
			})
		})

		Context("with Consumes", func() {
			const consumesMT = "application/json"

			BeforeEach(func() {
				dsl = func() {
					apidsl.Consumes("application/json")
				}
			})

			It("sets the API consumes", func() {
				Ω(Design.Consumes).Should(HaveLen(1))
				Ω(Design.Consumes[0].MIMETypes).Should(Equal([]string{consumesMT}))
				Ω(Design.Consumes[0].PackagePath).Should(BeEmpty())
			})

			Context("using a custom encoding package", func() {
				const pkgPath = "github.com/kyokomi/goa-v1/encoding/json"
				const fn = "NewFoo"

				BeforeEach(func() {
					dsl = func() {
						apidsl.Consumes("application/json", func() {
							apidsl.Package(pkgPath)
							apidsl.Function(fn)
						})
					}
				})

				It("sets the API consumes", func() {
					Ω(Design.Consumes).Should(HaveLen(1))
					Ω(Design.Consumes[0].MIMETypes).Should(Equal([]string{consumesMT}))
					Ω(Design.Consumes[0].PackagePath).Should(Equal(pkgPath))
					Ω(Design.Consumes[0].Function).Should(Equal(fn))
				})
			})
		})

		Context("with a BasePath", func() {
			const basePath = "basePath"

			BeforeEach(func() {
				dsl = func() {
					apidsl.BasePath(basePath)
				}
			})

			It("sets the API base path", func() {
				Ω(Design.BasePath).Should(Equal(basePath))
			})
		})

		Context("with Params", func() {
			const param1Name = "accountID"
			const param1Type = Integer
			const param1Desc = "the account ID"
			const param2Name = "id"
			const param2Type = String
			const param2Desc = "the widget ID"

			BeforeEach(func() {
				dsl = func() {
					apidsl.Params(func() {
						apidsl.Param(param1Name, param1Type, param1Desc)
						apidsl.Param(param2Name, param2Type, param2Desc)
					})
				}
			})

			It("sets the API base parameters", func() {
				Ω(Design.Params.Type).Should(BeAssignableToTypeOf(Object{}))
				params := Design.Params.Type.ToObject()
				Ω(params).Should(HaveLen(2))
				Ω(params).Should(HaveKey(param1Name))
				Ω(params).Should(HaveKey(param2Name))
				Ω(params[param1Name].Type).Should(Equal(param1Type))
				Ω(params[param2Name].Type).Should(Equal(param2Type))
				Ω(params[param1Name].Description).Should(Equal(param1Desc))
				Ω(params[param2Name].Description).Should(Equal(param2Desc))
			})

			Context("and a BasePath using them", func() {
				const basePath = "/:accountID/:id"

				BeforeEach(func() {
					prevDSL := dsl
					dsl = func() {
						apidsl.BasePath(basePath)
						prevDSL()
					}
				})

				It("sets both the base path and parameters", func() {
					Ω(Design.Params.Type).Should(BeAssignableToTypeOf(Object{}))
					params := Design.Params.Type.ToObject()
					Ω(params).Should(HaveLen(2))
					Ω(params).Should(HaveKey(param1Name))
					Ω(params).Should(HaveKey(param2Name))
					Ω(params[param1Name].Type).Should(Equal(param1Type))
					Ω(params[param2Name].Type).Should(Equal(param2Type))
					Ω(params[param1Name].Description).Should(Equal(param1Desc))
					Ω(params[param2Name].Description).Should(Equal(param2Desc))
					Ω(Design.BasePath).Should(Equal(basePath))
				})

				Context("with conflicting resource and API base params", func() {
					JustBeforeEach(func() {
						apidsl.Resource("foo", func() {
							apidsl.BasePath("/:accountID")
						})
						dslengine.Run()
					})

					It("returns an error", func() {
						Ω(dslengine.Errors).Should(HaveOccurred())
					})
				})

				Context("with an absolute resource base path", func() {
					JustBeforeEach(func() {
						apidsl.Resource("foo", func() {
							apidsl.Params(func() {
								apidsl.Param(param1Name, param1Type, param1Desc)
							})
							apidsl.BasePath("//:accountID")
						})
						dslengine.Run()
					})

					It("does not return an error", func() {
						Ω(dslengine.Errors).ShouldNot(HaveOccurred())
					})
				})
			})
		})

		Context("with ResponseTemplates", func() {
			const respName = "NotFound2"
			const respDesc = "Resource Not Found"
			const respStatus = 404
			const respMediaType = "text/plain"
			const respTName = "OK"
			const respTDesc = "All good"
			const respTStatus = 200

			BeforeEach(func() {
				dsl = func() {
					apidsl.ResponseTemplate(respName, func() {
						apidsl.Description(respDesc)
						apidsl.Status(respStatus)
						apidsl.Media(respMediaType)
					})
					apidsl.ResponseTemplate(respTName, func(mt string) {
						apidsl.Description(respTDesc)
						apidsl.Status(respTStatus)
						apidsl.Media(mt)
					})
				}
			})

			It("sets the API responses and response templates", func() {
				Ω(Design.Responses).Should(HaveKey(respName))
				Ω(Design.Responses[respName]).ShouldNot(BeNil())
				expected := ResponseDefinition{
					Name:        respName,
					Description: respDesc,
					Status:      respStatus,
					MediaType:   respMediaType,
				}
				actual := *Design.Responses[respName]
				Ω(actual).Should(Equal(expected))

				Ω(Design.ResponseTemplates).Should(HaveLen(1))
				Ω(Design.ResponseTemplates).Should(HaveKey(respTName))
				Ω(Design.ResponseTemplates[respTName]).ShouldNot(BeNil())
			})
		})

		Context("with Traits", func() {
			const traitName = "Authenticated"

			BeforeEach(func() {
				dsl = func() {
					apidsl.Trait(traitName, func() {
						apidsl.Headers(func() {
							apidsl.Header("Auth-Token")
							apidsl.Required("Auth-Token")
						})
					})
				}
			})

			It("sets the API traits", func() {
				Ω(Design.Traits).Should(HaveLen(1))
				Ω(Design.Traits).Should(HaveKey(traitName))
			})
		})

		Context("using Traits", func() {
			const traitName = "Authenticated"

			BeforeEach(func() {
				dsl = func() {
					apidsl.Trait(traitName, func() {
						apidsl.Attributes(func() {
							apidsl.Attribute("foo")
						})
					})
				}
			})

			JustBeforeEach(func() {
				apidsl.API(name, dsl)
				apidsl.MediaType("application/vnd.foo", func() {
					apidsl.UseTrait(traitName)
					apidsl.Attributes(func() {
						apidsl.Attribute("bar")
					})
					apidsl.View("default", func() {
						apidsl.Attribute("bar")
						apidsl.Attribute("foo")
					})
				})
				dslengine.Run()
			})

			It("sets the API traits", func() {
				Ω(Design.Traits).Should(HaveLen(1))
				Ω(Design.Traits).Should(HaveKey(traitName))
				Ω(Design.MediaTypes).Should(HaveKey("application/vnd.foo"))
				foo := Design.MediaTypes["application/vnd.foo"]
				Ω(foo.Type.ToObject()).ShouldNot(BeNil())
				o := foo.Type.ToObject()
				Ω(o).Should(HaveKey("foo"))
				Ω(o).Should(HaveKey("bar"))
			})
		})

		Context("using variadic Traits", func() {
			const traitName1 = "Authenticated"
			const traitName2 = "AuthenticatedTwo"

			BeforeEach(func() {
				dsl = func() {
					apidsl.Trait(traitName1, func() {
						apidsl.Attributes(func() {
							apidsl.Attribute("foo")
						})
					})
					apidsl.Trait(traitName2, func() {
						apidsl.Attributes(func() {
							apidsl.Attribute("baz")
						})
					})
				}
			})

			JustBeforeEach(func() {
				apidsl.API(name, dsl)
				apidsl.MediaType("application/vnd.foo", func() {
					apidsl.UseTrait(traitName1, traitName2)
					apidsl.Attributes(func() {
						apidsl.Attribute("bar")
					})
					apidsl.View("default", func() {
						apidsl.Attribute("bar")
						apidsl.Attribute("foo")
						apidsl.Attribute("baz")
					})
				})
				dslengine.Run()
			})

			It("sets the API traits", func() {
				Ω(Design.Traits).Should(HaveLen(2))
				Ω(Design.Traits).Should(HaveKey(traitName1))
				Ω(Design.Traits).Should(HaveKey(traitName2))
				Ω(Design.MediaTypes).Should(HaveKey("application/vnd.foo"))
				foo := Design.MediaTypes["application/vnd.foo"]
				Ω(foo.Type.ToObject()).ShouldNot(BeNil())
				o := foo.Type.ToObject()
				Ω(o).Should(HaveKey("foo"))
				Ω(o).Should(HaveKey("bar"))
				Ω(o).Should(HaveKey("baz"))
			})
		})
	})

})
