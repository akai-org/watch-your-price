import requests
from django.http import HttpResponseBadRequest
from rest_framework import status
from rest_framework.response import Response
from rest_framework.views import APIView


class SearchProducts(APIView):

    def get(self, request):

        params = {}

        try:
            params['phrase'] = request.query_params['phrase']
            params['website'] = request.query_params['website']
        except KeyError:
            return HttpResponseBadRequest()
        params['page'] = request.query_params['page'] if 'page' in request.query_params else 0

        try:
            response = requests.get('http://localhost:8001/search', params=params)
            response.raise_for_status()
        except requests.exceptions.RequestException:
            return HttpResponseBadRequest()

        return Response(response.json(), status=status.HTTP_200_OK)
