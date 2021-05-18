from django.http import HttpResponseBadRequest
from rest_framework import status
from rest_framework.response import Response
from rest_framework.views import APIView


class SearchProducts(APIView):

    def get(self, request):
        try:
            phrase = request.query_params['phrase']
            website = request.query_params['website']
        except KeyError:
            return HttpResponseBadRequest()
        page = request.query_params['page'] if 'page' in request.query_params else 0
        return Response(request.query_params, status=status.HTTP_200_OK)
